
#ifndef OPERATION_H
#define OPERATION_H

#include <QWidget>
#include<QString>
#include<listuser.h>
#include<QUdpSocket>
#include <QtWebSockets/QWebSocket>
#include<QFile>
#include<QCryptographicHash>
#include<QDebug>
#include<QTextStream>
#include<QDir>
#include<QResource>
#include <QRandomGenerator>
#include<QList>
#include <QNetworkInterface>
#include<QHostInfo>

QT_BEGIN_NAMESPACE
namespace Ui { class Operation; }
QT_END_NAMESPACE

class Operation : public QWidget

{
    Q_OBJECT

public:
    Operation(QWidget *parent = nullptr);
    //建立本地websocket连接
    void LocalConnect();
    //哈希加密
    void SaveEncryptData(const QString &data);
    QString getLocalIPAddress();
    ~Operation();

private:
    Ui::Operation *ui;
    ListUser *list_user;
    //本地通信
    QWebSocket *webSocket;
    //身份验证
    QString username;
    QString identity;
};

#endif // OPERATION_H
