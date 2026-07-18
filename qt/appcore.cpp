#include "appcore.h"
#include "ipcclient.h"
#include <QRandomGenerator>

AppCore::AppCore(QObject *parent)
    : QObject(parent)
    , m_ipc(new IpcClient(this))
    , m_isConnected(false)
    , m_statusText("Tap to Connect")
    , m_selectedServer("US - New York #1")
    , m_mockRx(0.0)
    , m_mockTx(0.0)
    , m_mockPing(45)
{
    connect(m_ipc, &IpcClient::commandSuccess, this, [this](const QString &msg) {
        if (m_isConnected) {
            onDisconnectSuccess(msg);
        } else {
            onConnectSuccess(msg);
        }
    });
    connect(m_ipc, &IpcClient::commandError, this, &AppCore::onCommandError);

    // Real-time Stats Timer (Updates every 1 second like a real VPN)
    m_statsTimer = new QTimer(this);
    connect(m_statsTimer, &QTimer::timeout, this, &AppCore::updateMockStats);
    m_statsTimer->start(1000);
}

bool AppCore::isConnected() const { return m_isConnected; }
QString AppCore::statusText() const { return m_statusText; }
QString AppCore::selectedServer() const { return m_selectedServer; }
QString AppCore::rxData() const { return formatBytes(m_mockRx); }
QString AppCore::txData() const { return formatBytes(m_mockTx); }
QString AppCore::pingData() const { return QString::number(m_mockPing) + " ms"; }

void AppCore::setSelectedServer(const QString &server) {
    if (m_selectedServer != server) {
        m_selectedServer = server;
        emit selectedServerChanged();
        // Simulate ping change on server switch
        m_mockPing = QRandomGenerator::global()->bounded(20, 120);
        emit statsUpdated();
    }
}

void AppCore::toggleConnection() {
    if (m_isConnected) {
        m_ipc->sendCommand("disconnect");
    } else {
        m_ipc->sendCommand("connect", m_selectedServer);
    }
}

void AppCore::openServerList() {
    emit showServerListRequested();
}

void AppCore::onConnectSuccess(const QString & /*message*/) {
    m_isConnected = true;
    m_statusText = "Connected";
    emit connectionStateChanged();
}

void AppCore::onDisconnectSuccess(const QString & /*message*/) {
    m_isConnected = false;
    m_statusText = "Tap to Connect";
    emit connectionStateChanged();
}

void AppCore::onCommandError(const QString & /*errorMessage*/) {
    m_statusText = "Error";
    emit connectionStateChanged();
}

void AppCore::updateMockStats() {
    if (m_isConnected) {
        // Simulate realistic VPN traffic (MB/s)
        m_mockRx += (QRandomGenerator::global()->bounded(100000, 500000) / 1024.0 / 1024.0); 
        m_mockTx += (QRandomGenerator::global()->bounded(50000, 200000) / 1024.0 / 1024.0);
        // Ping fluctuation
        if(QRandomGenerator::global()->bounded(0, 10) == 0) {
            m_mockPing = QRandomGenerator::global()->bounded(20, 100);
        }
    }
    emit statsUpdated();
}

QString AppCore::formatBytes(double bytes) const {
    if (bytes < 1024) return QString::number(bytes, 'f', 2) + " B";
    if (bytes < 1048576) return QString::number(bytes / 1024.0, 'f', 2) + " MB";
    if (bytes < 1073741824) return QString::number(bytes / 1048576.0, 'f', 2) + " MB";
    return QString::number(bytes / 1073741824.0, 'f', 2) + " GB";
}
