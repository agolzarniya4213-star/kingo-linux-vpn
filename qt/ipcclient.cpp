#include "ipcclient.h"
#include <QJsonDocument>
#include <QJsonObject>
#include <QTimer>
#include <QQueue>

IpcClient::IpcClient(QObject *parent) : QObject(parent), m_socket(new QLocalSocket(this)) {
    connect(m_socket, &QLocalSocket::readyRead, this, &IpcClient::onReadyRead);
    connect(m_socket, &QLocalSocket::connected, this, &IpcClient::onConnected);
    connect(m_socket, &QLocalSocket::disconnected, this, [this]() {
        // تلاش مجدد برای اتصال پس از 3 ثانیه
        QTimer::singleShot(3000, this, [this]() {
            if (m_socket->state() != QLocalSocket::ConnectedState && m_socket->state() != QLocalSocket::ConnectingState) {
                m_socket->connectToServer("/tmp/kingo-vpn.sock");
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
        // صف کردن درخواست‌ها برای جلوگیری از Data Loss
        m_pendingRequests.enqueue(data);
        if (m_socket->state() != QLocalSocket::ConnectingState) {
            m_socket->connectToServer("/tmp/kingo-vpn.sock");
        }
    }
}

void IpcClient::onConnected() {
    // ارسال تمام درخواست‌های معلق در صف
    while (!m_pendingRequests.isEmpty()) {
        m_socket->write(m_pendingRequests.dequeue());
    }
    m_socket->flush();
}

void IpcClient::onReadyRead() {
    m_buffer.append(m_socket->readAll());
    // محدودیت حجم Buffer برای جلوگیری از OOM
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
