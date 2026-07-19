#include "ipcclient.h"
#include <QJsonDocument>
#include <QJsonObject>
#include <QTimer>

IpcClient::IpcClient(QObject *parent) : QObject(parent), m_socket(new QLocalSocket(this)) {
    connect(m_socket, &QLocalSocket::readyRead, this, &IpcClient::onReadyRead);
    connect(m_socket, &QLocalSocket::connected, this, &IpcClient::onConnected);
    connect(m_socket, &QLocalSocket::disconnected, this, [this]() {
        QTimer::singleShot(3000, this, [this]() {
            if (m_socket->state() != QLocalSocket::ConnectedState && m_socket->state() != QLocalSocket::ConnectingState) {
                // FIX: Point to the new secure socket path
                m_socket->connectToServer("/run/kingo-vpn/kingo-vpn.sock");
            }
        });
    });
    connect(m_socket, &QLocalSocket::errorOccurred, this, [this](QLocalSocket::LocalSocketError socketError) {
        Q_UNUSED(socketError);
        emit errorOccurred(m_socket->errorString());
    });
}

void IpcClient::sendRequest(const QJsonObject &request) {
    QByteArray data = QJsonDocument(request).toJson(QJsonDocument::Compact) + "\n";
    
    if (m_socket->state() == QLocalSocket::ConnectedState) {
        m_socket->write(data);
        m_socket->flush();
    } else {
        m_pendingRequests.enqueue(data);
        if (m_socket->state() != QLocalSocket::ConnectingState) {
            m_socket->connectToServer("/run/kingo-vpn/kingo-vpn.sock");
        }
    }
}

void IpcClient::onConnected() {
    while (!m_pendingRequests.isEmpty()) {
        m_socket->write(m_pendingRequests.dequeue());
    }
    m_socket->flush();
}

void IpcClient::onReadyRead() {
    m_buffer.append(m_socket->readAll());
    if (m_buffer.size() > 1024 * 1024) { 
        m_buffer.clear();
        emit errorOccurred("IPC Buffer overflow detected!");
        return;
    }
    
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
