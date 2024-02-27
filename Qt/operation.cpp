
#include "operation.h"
#include "ui_operation.h"
#include<QMessageBox>

Operation::Operation(QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::Operation)
{
    ui->setupUi(this);
    //建立本地连接
    LocalConnect();
    //随机数生成
    QRandomGenerator::securelySeeded();
    //局域网内部ip地址基本上不相同，利用这个作为种子生成随机数
    QString WLAN_IPV4 = getLocalIPAddress();
    qDebug()<<WLAN_IPV4;
    QByteArray byteArray = WLAN_IPV4.toUtf8();
    QByteArray hash = QCryptographicHash::hash(byteArray, QCryptographicHash::Sha256);
    quint32 uniqueValue = qHash(hash);
    qDebug()<<uniqueValue;
    //使用固定的种子，使得在同一局域网的同一台设备下ip地址相同
    QRandomGenerator randomGenerator(uniqueValue);
    int randomNumber = randomGenerator.bounded(1000000000);
    identity = QString(randomNumber);
    connect(ui->EnterButton,&QPushButton::clicked,this,[=]{
        QString userName = ui->userName->text();
        QFile file("username.txt");
        //可读写方式打开文件
        if(file.open(QIODevice::ReadWrite | QIODevice::Text)){
            QTextStream in(&file);
            username = in.readAll().trimmed(); // 读取文件中的文本并去除首尾空格
            if(in.readAll().isEmpty())
            {
                in.seek(0);
                in<<userName;
            }
            else
            {
                //界面切换
                if(userName == username){
                    this->hide();
                    this->list_user = new ListUser();
                    this->list_user->show();
                }
            }
            file.close();
        }
        else{
            QMessageBox::warning(this,"提示","username.txt文件无法打开");
            return;
        }
    });
    //将指纹标识到文件中
    SaveEncryptData(QString(randomNumber));
}

void Operation::LocalConnect()
{
    webSocket = new QWebSocket();
    QUrl url("ws://192.168.203.179:80/ws");
    webSocket->open(url);
}

//将输入加密后保存到本地
//注意一点，资源文件是只读的，不能直接写入
void Operation::SaveEncryptData(const QString &data)
{
    QByteArray byteArray = data.toUtf8();
    QByteArray hash = QCryptographicHash::hash(byteArray, QCryptographicHash::Sha256);
    QFile file("identity.txt");
    if (file.open(QIODevice::WriteOnly | QIODevice::Text)) {
        // 创建 QTextStream 对象，用于写入文件
        QTextStream out(&file);
        // 将哈希值转换为十六进制字符串，并写入文件
        out << hash.toHex();
        file.close();
        qDebug() << "已成功保存到 identity.txt";
    } else {
        qDebug() << "无法打开" << file.errorString();
    }
}

//获取IPV4地址
QString Operation::getLocalIPAddress()
{
    QList<QNetworkInterface> interfaces = QNetworkInterface::allInterfaces();
    qDebug() << "";
    for (const QNetworkInterface &interface : interfaces) {
        if (interface.name().startsWith("wireless", Qt::CaseInsensitive)) {
            QList<QNetworkAddressEntry> entries = interface.addressEntries();
            for (const QNetworkAddressEntry &entry : entries) {
                if (entry.ip().protocol() == QAbstractSocket::IPv4Protocol) {
                    QString ipv4Address = entry.ip().toString();
                    if (ipv4Address.startsWith("192.168")) {
                        //qDebug() << "" << interface.name();
                        //qDebug() << "" << ipv4Address;
                        return ipv4Address;
                    }
                }
            }
        }
    }


    return QString();
}

Operation::~Operation()
{
    delete ui;
}


