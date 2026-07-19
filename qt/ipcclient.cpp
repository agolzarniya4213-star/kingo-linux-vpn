#include "ipcclient.h"
#include <QJsonDocument>
#include <QJsonObject>

IpcClient::IpcClient(QObject *parent) : QObject(parent), m_socket(new QLocalSocket(this)) {
    connect(m_socket, &QLocalSocket::readyRead, this, &IpcClient::onReadyRead);
    connect(m_socket, &QLocalSocket::connected, this, &IpcClient::onConnected);
    connect(m_socket, &QLocalSocket::errorOccurred, this, [this](QLocalSocket::LocalSocketError socketError) {
        Q_UNUSED(socketError);
        emit errorOccurred(m_socket->errorString());
    });
}

void IpcClient::sendRequest(const QJsonObject &request) {
    if (m_socket->state() != QLocalSocket::ConnectedState) {
        m_pendingRequest = QJsonDocument(request).toJson(QJsonDocument::Compact) + "\n";
        m_socket->connectToServer("/tmp/kingo-vpn.sock");
        return;
    }
    
    QByteArray data = QJsonDocument(request).toJson(QJsonDocument::Compact) + "\n";
    m_socket->write(data);
    m_socket->flush();
}

void IpcClient::onConnected() {
    if (!m_pendingRequest.isEmpty()) {
        m_socket->write(m_pendingRequest);
        m_socket->flush();
        m_pendingRequest.clear();
    }
}

void IpcClient::onReadyRead() {
    m_buffer.append(m_socket->readAll());
    int newlineIndex;
    while ((newlineIndex = m_buffer.indexOf('\n')) != -1) {
        QByteArray line = m_buffer.left(newlineIndex);
        m_buffer.remove(0, newlineIndex + 1);
        
        QJsonParseError err;
        QJsonDocument doc = QJsonDocument::fromJson(line, &err);
        if (err.error == QJsonParseError::NoError && doc.isObject()) {
            emit responseReceived(doc.object());
        }
    }
}
