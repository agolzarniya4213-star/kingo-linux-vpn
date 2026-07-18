#include "appcore.h"
#include "ipcclient.h"

AppCore::AppCore(QObject *parent)
    : QObject(parent)
    , m_ipc(new IpcClient(this))
    , m_timer(new QTimer(this))
    , m_connectionStatus("Disconnected")
    , m_downloadSpeed("0 B/s")
    , m_uploadSpeed("0 B/s")
    , m_totalDownload("0 B")
    , m_totalUpload("0 B")
{
    connect(m_timer, &QTimer::timeout, this, &AppCore::pollStatus);
    m_timer->start(1000);
}

AppCore::~AppCore() = default;

QString AppCore::connectionStatus() const { return m_connectionStatus; }
QString AppCore::downloadSpeed() const { return m_downloadSpeed; }
QString AppCore::uploadSpeed() const { return m_uploadSpeed; }
QString AppCore::totalDownload() const { return m_totalDownload; }
QString AppCore::totalUpload() const { return m_totalUpload; }
QString AppCore::selectedServer() const { return m_selectedServer; }
QVariantList AppCore::serverList() const { return m_serverList; }

void AppCore::setSelectedServer(const QString &server) {
    if (m_selectedServer != server) {
        m_selectedServer = server;
        emit selectedServerChanged();
    }
}

void AppCore::connectVPN() {
    m_ipc->sendCommand("connect");
    m_connectionStatus = "Connecting";
    emit statusChanged();
}

void AppCore::disconnectVPN() {
    m_ipc->sendCommand("disconnect");
    m_connectionStatus = "Disconnected";
    emit statusChanged();
}

void AppCore::refreshServers() {
    m_ipc->sendCommand("refresh");
}

QString AppCore::formatBytes(double bytes) const {
    if (bytes < 0.0) { bytes = 0.0; }
    static const char *units[] = {"B", "KB", "MB", "GB", "TB"};
    int index = 0;
    double size = bytes;
    while (size >= 1024.0 && index < 4) {
        size /= 1024.0;
        ++index;
    }
    return QString("%1 %2").arg(size, 0, 'f', 1).arg(units[index]);
}

void AppCore::pollStatus() {
    m_ipc->sendCommand("status");
}
