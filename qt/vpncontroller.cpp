#include "vpncontroller.h"
#include <QDateTime>
#include <QRandomGenerator>

VpnController::VpnController(QObject *parent) : QObject(parent), m_client(new IpcClient(this)) {
    connect(m_client, &IpcClient::responseReceived, this, &VpnController::onResponseReceived);
    connect(m_client, &IpcClient::errorOccurred, this, [this](const QString &err){
        appendLog("IPC Error: " + err);
    });
    refreshStatus();
    fetchServers();
}

void VpnController::appendLog(const QString &log) {
    m_logs += QDateTime::currentDateTime().toString("hh:mm:ss") + " - " + log + "\n";
    emit logsChanged();
}

QString VpnController::generateRequestID() {
    return QString::number(QDateTime::currentMSecsSinceEpoch()) + "-" + QString::number(QRandomGenerator::global()->generate());
}

void VpnController::connectToServer(const QString &uri) {
    appendLog("Connecting to server...");
    QJsonObject req; req["request_id"] = generateRequestID(); req["action"] = "connect_server"; req["server_uri"] = uri;
    m_client->sendRequest(req);
}

void VpnController::autoConnect() {
    appendLog("Finding best server...");
    QJsonObject req; req["request_id"] = generateRequestID(); req["action"] = "auto_connect";
    m_client->sendRequest(req);
}

void VpnController::disconnectVpn() {
    appendLog("Disconnecting...");
    QJsonObject req; req["request_id"] = generateRequestID(); req["action"] = "disconnect";
    m_client->sendRequest(req);
}

void VpnController::refreshStatus() {
    QJsonObject req; req["request_id"] = generateRequestID(); req["action"] = "status";
    m_client->sendRequest(req);
}

void VpnController::fetchServers() {
    QJsonObject req; req["request_id"] = generateRequestID(); req["action"] = "get_servers";
    m_client->sendRequest(req);
}

void VpnController::addSubscription(const QString &url) {
    appendLog("Updating subscription...");
    QJsonObject req; req["request_id"] = generateRequestID(); req["action"] = "add_subscription"; req["sub_url"] = url;
    m_client->sendRequest(req);
}

void VpnController::clearServers() {
    appendLog("Clearing servers...");
    // FIX: Clear local list immediately for instant UI feedback
    m_servers.clear();
    emit serversChanged();
    
    QJsonObject req; req["request_id"] = generateRequestID(); req["action"] = "clear_servers";
    m_client->sendRequest(req);
}

void VpnController::testLatency() {
    appendLog("Testing latency...");
    QJsonObject req; req["request_id"] = generateRequestID(); req["action"] = "test_latency";
    m_client->sendRequest(req);
}

void VpnController::getTraffic() {
    QJsonObject req; req["request_id"] = generateRequestID(); req["action"] = "get_traffic";
    m_client->sendRequest(req);
}

void VpnController::onResponseReceived(const QJsonObject &response) {
    if (response.contains("state")) {
        setStatus(response["state"].toString());
    }
    if (response.contains("servers")) {
        m_servers = response["servers"].toVariant().toList();
        emit serversChanged();
        appendLog("Server list updated (" + QString::number(m_servers.size()) + " items).");
    }
    if (response.contains("upload") && response.contains("download")) {
        m_uploadSpeed = response["upload"].toVariant().toLongLong();
        m_downloadSpeed = response["download"].toVariant().toLongLong();
        emit trafficChanged();
    }
    if (response.contains("success") && !response["success"].toBool()) {
        if (response.contains("message") && !response["message"].toString().isEmpty()) {
            appendLog("Error: " + response["message"].toString());
            emit errorOccurred(response["message"].toString());
        }
    }
}

void VpnController::setStatus(const QString &newStatus) {
    if (m_status != newStatus) {
        m_status = newStatus;
        emit statusChanged();
        appendLog("Status changed to: " + newStatus.toUpper());
    }
}
