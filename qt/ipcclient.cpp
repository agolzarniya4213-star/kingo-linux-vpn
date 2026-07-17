#include "ipcclient.h"
#include <QStandardPaths>
#include <QJsonDocument>
#include <QDir>

IpcClient::IpcClient(QObject *parent)
    : QObject(parent)
    , m_socket(new QLocalSocket(this))
{
    initializeSocketPath();
    connect(m_socket, &QLocalSocket::errorOccurred, 
            this, &IpcClient::onSocketError);
}

void IpcClient::initializeSocketPath()
{
    QString configDir = QStandardPaths::writableLocation(QStandardPaths::GenericConfigLocation);
    m_socketPath = configDir + QDir::separator() + "kingo-linux-vpn" + QDir::separator() + "daemon.sock";
}

void IpcClient::sendCommand(const QString &action)
{
    if (action.isEmpty()) {
        emit commandError("Internal Error: Action cannot be empty.");
        return;
    }

    m_socket->connectToServer(m_socketPath);

    if (!m_socket->waitForConnected(500)) {
        emit commandError("Cannot connect to Kingo daemon. Is it running?");
        return;
    }

    QJsonObject payload;
    payload["action"] = action;
    QJsonDocument doc(payload);
    QByteArray data = doc.toJson(QJsonDocument::Compact) + "\n";

    m_socket->write(data);
    m_socket->flush();

    if (!m_socket->waitForReadyRead(1000)) {
        emit commandError("Daemon timed out.");
        m_socket->disconnectFromServer();
        return;
    }

    QByteArray response = m_socket->readAll();
    m_socket->disconnectFromServer();

    QJsonObject resObj = parseResponse(response);
    
    if (resObj["status"].toString() == "success") {
        emit commandSuccess(resObj["message"].toString());
    } else {
        emit commandError(resObj["message"].toString("Unknown daemon error."));
    }
}

void IpcClient::onSocketError(QLocalSocket::LocalSocketError error)
{
    Q_UNUSED(error)
    emit commandError(m_socket->errorString());
}

QJsonObject IpcClient::parseResponse(const QByteArray &data)
{
    QJsonDocument doc = QJsonDocument::fromJson(data);
    if (!doc.isObject()) {
        return {{"status", "error"}, {"message", "Invalid JSON from daemon."}};
    }
    return doc.object();
}
