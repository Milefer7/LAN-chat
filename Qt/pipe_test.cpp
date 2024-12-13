#include <QCoreApplication>
#include <QLocalSocket>
#include <QDebug>
#include <QDateTime>

int main(int argc, char *argv[])
{
    QCoreApplication a(argc, argv);

    QLocalSocket socket;
    socket.connectToServer("\\\\.\\pipe\\my_pipe");

    if (socket.waitForConnected()) {
        qDebug() << "Connected!";

        while (socket.waitForReadyRead()) {
            QByteArray byteArray = socket.readAll();
            QList<QByteArray> messageAndTimestamp = byteArray.split('|');

            QString message = QString::fromUtf8(messageAndTimestamp[0]);
            qint64 timestamp = messageAndTimestamp[1].toLongLong();

            qDebug() << "Received message:" << message;
            qDebug() << "Received timestamp:" << timestamp;

            qint64 currentTimestamp = QDateTime::currentMSecsSinceEpoch();
            qDebug() << "Current timestamp:" << currentTimestamp;

            qDebug() << "Elapsed time:" << currentTimestamp - timestamp << "ms";
        }
    } else {
        qDebug() << "Failed to connect:" << socket.errorString();
    }

    return a.exec();
}
