#ifndef LOCALSERVER_H
#define LOCALSERVER_H

#include <QWidget>
#include<QWebSocketServer>
#include<QWebSocket>
#include<QList>
#include<QMessageBox>
#include<QString>

namespace Ui {
class LocalServer;
}

class LocalServer : public QWidget
{
    Q_OBJECT

public:
    explicit LocalServer(QWidget *parent = nullptr);
    ~LocalServer();

private:
    Ui::LocalServer *ui;
    QWebSocketServer *server;
    QList<QWebSocket *> clients;

private slots:
    void onNewConnection();
    void processMessage(const QString &message);
    void processBinaryMessage(const QByteArray &message);
    void onSocketDisconnected();


};

#endif // LOCALSERVER_H
