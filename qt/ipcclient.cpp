#include "ipcclient.h"
#include <QJsonDocument>
#include <QJsonObject>
#include <QStandardPaths>

IpcClient::IpcClient(QObject *parent)
    : QObject(parent)
    , m_socket(new QLocalSocket(this))
{
    // Find the exact same path as the Go daemon
    QString configDir = QStandardPaths::writableLocation(QStandardPaths::GenericConfigLocation);
    m_socketPath = configDir + "/kingo-linux-vpn/daemon.sock";
}

void IpcClient::sendCommand(const QString &action)
{
    m_socket->connectToServer(m_socketPath);
    
    // Wait 500ms for the daemon to respond
    if (!m_socket->waitForConnected(500)) {
        emit commandError("Cannot connect to daemon. Is it running?");
        return;
    }

    // Create JSON payload: {"action": "connect"}
    QJsonObject obj;
    obj["action"] = action;
    QJsonDocument doc(obj);
    QByteArray data = doc.toJson(QJsonDocument::Compact) + "\n";

    m_socket->write(data);
    m_socket->flush();

    if (!m_socket->waitForReadyRead(1000)) {
        emit commandError("Daemon timeout");
        m_socket->disconnectFromServer();
        return;
    }

    // Read response
    QByteArray response = m_socket->readAll();
    m_socket->disconnectFromServer();

    QJsonDocument resDoc = QJsonDocument::fromJson(response);
    if (resDoc.object()["status"].toString() == "success") {
        emit commandSuccess(resDoc.object()["message"].toString());
    } else {
        emit commandError(resDoc.object()["message"].toString());
    }
}
