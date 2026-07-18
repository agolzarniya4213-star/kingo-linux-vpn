#include "ipcclient.h"
#include <QJsonDocument>
#include <QJsonObject>

IpcClient::IpcClient(QObject *parent) : QObject(parent), m_socket(new QLocalSocket(this)) {
    connect(m_socket, &QLocalSocket::readyRead, this, &IpcClient::onReadyRead);
    connect(m_socket, &QLocalSocket::errorOccurred, this, [this](QLocalSocket::LocalSocketError socketError) {
        Q_UNUSED(socketError);
        emit errorOccurred(m_socket->errorString());
    });
}

void IpcClient::sendRequest(const QJsonObject &request) {
    if (m_socket->state() != QLocalSocket::ConnectedState) {
        m_socket->connectToServer("/tmp/kingo-vpn.sock");
    }
    
    if (m_socket->state() == QLocalSocket::ConnectingState || m_socket->state() == QLocalSocket::ConnectedState) {
        QJsonDocument doc(request);
        m_socket->write(doc.toJson(QJsonDocument::Compact) + "\n");
        m_socket->flush();
    } else {
        emit errorOccurred("Failed to connect to daemon.");
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
