#include "ipcclient.h"
#include <QJsonDocument>
#include <QJsonObject>
#include <QTimer>

IpcClient::IpcClient(QObject *parent) : QObject(parent), m_socket(new QLocalSocket(this)), m_useTmpPath(false) {
    connect(m_socket, &QLocalSocket::readyRead, this, &IpcClient::onReadyRead);
    connect(m_socket, &QLocalSocket::connected, this, &IpcClient::onConnected);
    connect(m_socket, &QLocalSocket::disconnected, this, [this]() {
        QTimer::singleShot(3000, this, [this]() {
            attemptConnection();
        });
    });
    connect(m_socket, &QLocalSocket::errorOccurred, this, [this](QLocalSocket::LocalSocketError socketError) {
        Q_UNUSED(socketError);
        // If /run fails, fallback to /tmp on next attempt
        m_useTmpPath = !m_useTmpPath;
        QTimer::singleShot(1000, this, [this]() {
            attemptConnection();
        });
    });
}

void IpcClient::attemptConnection() {
    if (m_socket->state() != QLocalSocket::ConnectedState && m_socket->state() != QLocalSocket::ConnectingState) {
        QString path = m_useTmpPath ? "/tmp/kingo-vpn/kingo-vpn.sock" : "/run/kingo-vpn/kingo-vpn.sock";
        m_socket->connectToServer(path);
    }
}

void IpcClient::sendRequest(const QJsonObject &request) {
    QByteArray data = QJsonDocument(request).toJson(QJsonDocument::Compact) + "\n";
    
    if (m_socket->state() == QLocalSocket::ConnectedState) {
        m_socket->write(data);
        m_socket->flush();
    } else {
        if (m_pendingRequests.size() < 100) {
            m_pendingRequests.enqueue(data);
        }
        attemptConnection();
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
