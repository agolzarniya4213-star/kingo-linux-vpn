#include "ipcclient.h"
#include <QJsonDocument>
#include <QJsonObject>
#include <QStandardPaths>
#include <QDir>

IpcClient::IpcClient(QObject *parent)
    : QObject(parent)
    , m_socket(new QLocalSocket(this))
{
    // Calculate exact socket path matching Go daemon logic
    QString configDir = QStandardPaths::writableLocation(QStandardPaths::GenericConfigLocation);
    m_socketPath = configDir + QDir::separator() + "kingo-linux-vpn" + QDir::separator() + "daemon.sock";

    // Async event-driven connections
    connect(m_socket, &QLocalSocket::connected, this, &IpcClient::onConnected);
    connect(m_socket, &QLocalSocket::readyRead, this, &IpcClient::onReadyRead);
    connect(m_socket, &QLocalSocket::errorOccurred, this, &IpcClient::onError);
}

void IpcClient::sendCommand(const QString &action)
{
    if (m_socket->state() == QLocalSocket::UnconnectedState) {
        m_pendingAction = action;
        m_socket->connectToServer(m_socketPath);
    } else {
        emit commandError("Socket is busy. Please wait.");
    }
}

void IpcClient::onConnected()
{
    // Build JSON payload
    QJsonObject obj;
    obj["action"] = m_pendingAction;
    QJsonDocument doc(obj);
    QByteArray data = doc.toJson(QJsonDocument::Compact) + "\n";

    m_socket->write(data);
    m_socket->flush();
}

void IpcClient::onReadyRead()
{
    QByteArray response = m_socket->readAll();
    m_socket->disconnectFromServer();

    QJsonDocument resDoc = QJsonDocument::fromJson(response);
    if (resDoc.object()["status"].toString() == "success") {
        emit commandSuccess(resDoc.object()["message"].toString());
    } else {
        emit commandError(resDoc.object()["message"].toString());
    }
}

void IpcClient::onError(QLocalSocket::LocalSocketError socketError)
{
    Q_UNUSED(socketError)
    m_socket->disconnectFromServer();
    emit commandError("Cannot connect to daemon. Is kingo-linux-vpn daemon running?");
}
