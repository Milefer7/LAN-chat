#include "localserver.h"
#include "ui_localserver.h"

LocalServer::LocalServer(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::LocalServer)
{
    ui->setupUi(this);
    //(描述性字符串,非安全WebSocket服务器，不适用加密协议,)
    server = new QWebSocketServer("Local WebSocket Server", QWebSocketServer::NonSecureMode, this);
    //server监听外部端口连接
    if (!server->listen(QHostAddress::Any, 12345)) { // 监听端口号为 12345
        qFatal("Failed to start WebSocket server: %s", qPrintable(server->errorString()));
        return;
    }
    connect(server,&QWebSocketServer::newConnection,this,&LocalServer::onNewConnection);
}


LocalServer::~LocalServer()
{
    delete ui;
}

//处理新连接
void LocalServer::onNewConnection()
{
    QWebSocket *clientSocket = server->nextPendingConnection();
    connect(clientSocket, &QWebSocket::textMessageReceived, this, &LocalServer::processMessage);
    connect(clientSocket, &QWebSocket::binaryMessageReceived, this, &LocalServer::processBinaryMessage);
    connect(clientSocket, &QWebSocket::disconnected, this, &LocalServer::onSocketDisconnected);
    clients.append(clientSocket);
}

void LocalServer::processMessage(const QString &message)
{

}

void LocalServer::processBinaryMessage(const QByteArray &message)
{

}

void LocalServer::onSocketDisconnected()
{

}


