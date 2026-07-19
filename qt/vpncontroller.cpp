#include "vpncontroller.h"
#include <QDateTime>
#include <QRandomGenerator>

VpnController::VpnController(QObject *parent) : QObject(parent), m_client(new IpcClient(this)) {
    connect(m_client, &IpcClient::responseReceived, this, &VpnController::onResponseReceived);
    connect(m_client, &IpcClient::errorOccurred, this, &VpnController::errorOccurred);
    refreshStatus();
    fetchServers();
}

QString VpnController::generateRequestID() {
    return QString::number(QDateTime::currentMSecsSinceEpoch()) + "-" + QString::number(QRandomGenerator::global()->generate());
}

void VpnController::connectToServer(const QString &uri) {
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "connect_server";
    req["server_uri"] = uri;
    m_client->sendRequest(req);
}

void VpnController::autoConnect() {
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "auto_connect";
    m_client->sendRequest(req);
}

void VpnController::disconnectVpn() {
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "disconnect";
    m_client->sendRequest(req);
}

void VpnController::refreshStatus() {
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "status";
    m_client->sendRequest(req);
}

void VpnController::fetchServers() {
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "get_servers";
    m_client->sendRequest(req);
}

void VpnController::addSubscription(const QString &url) {
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "add_subscription";
    req["sub_url"] = url;
    m_client->sendRequest(req);
}

void VpnController::testLatency() {
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "test_latency";
    m_client->sendRequest(req);
}

void VpnController::getTraffic() {
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "get_traffic";
    m_client->sendRequest(req);
}

void VpnController::onResponseReceived(const QJsonObject &response) {
    // Note: In a fully async UI, we would map response["request_id"] to a callback.
    // Since QML updates are atomic based on state, we process directly.
    
    if (response.contains("state")) {
        setStatus(response["state"].toString());
    }
    if (response.contains("servers")) {
        m_servers = response["servers"].toVariant().toList();
        emit serversChanged();
    }
    if (response.contains("upload") && response.contains("download")) {
        m_uploadSpeed = response["upload"].toVariant().toLongLong();
        m_downloadSpeed = response["download"].toVariant().toLongLong();
        emit trafficChanged();
    }
    if (response.contains("success") && !response["success"].toBool()) {
        if (response.contains("message") && !response["message"].toString().isEmpty()) {
            emit errorOccurred(response["message"].toString());
        }
    }
}

void VpnController::setStatus(const QString &newStatus) {
    if (m_status != newStatus) {
        m_status = newStatus;
        emit statusChanged();
    }
}
