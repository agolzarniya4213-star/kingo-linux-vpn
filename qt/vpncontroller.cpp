#include "vpncontroller.h"

VpnController::VpnController(QObject *parent) : QObject(parent), m_client(new IpcClient(this)) {
    connect(m_client, &IpcClient::responseReceived, this, &VpnController::onResponseReceived);
    connect(m_client, &IpcClient::errorOccurred, this, &VpnController::errorOccurred);
    refreshStatus();
    fetchServers();
}

void VpnController::connectToServer(const QString &uri) {
    QJsonObject req;
    req["action"] = "connect_server";
    req["server_uri"] = uri;
    m_client->sendRequest(req);
}

void VpnController::autoConnect() {
    QJsonObject req;
    req["action"] = "auto_connect";
    m_client->sendRequest(req);
}

void VpnController::disconnectVpn() {
    QJsonObject req;
    req["action"] = "disconnect";
    m_client->sendRequest(req);
}

void VpnController::refreshStatus() {
    QJsonObject req;
    req["action"] = "status";
    m_client->sendRequest(req);
}

void VpnController::fetchServers() {
    QJsonObject req;
    req["action"] = "get_servers";
    m_client->sendRequest(req);
}

void VpnController::addSubscription(const QString &url) {
    QJsonObject req;
    req["action"] = "add_subscription";
    req["sub_url"] = url;
    m_client->sendRequest(req);
}

void VpnController::testLatency() {
    QJsonObject req;
    req["action"] = "test_latency";
    m_client->sendRequest(req);
}

void VpnController::getTraffic() {
    QJsonObject req;
    req["action"] = "get_traffic";
    m_client->sendRequest(req);
}

void VpnController::onResponseReceived(const QJsonObject &response) {
    if (response.contains("state")) {
        setStatus(response["state"].toString());
    }
    if (response.contains("servers")) {
        m_servers = response["servers"].toVariant().toList();
        emit serversChanged();
    }
    if (response.contains("upload") || response.contains("download")) {
        m_uploadSpeed = response["upload"].toVariant().toLongLong();
        m_downloadSpeed = response["download"].toVariant().toLongLong();
        emit trafficChanged();
    }
    if (response.contains("message") && !response["message"].toString().isEmpty()) {
        emit errorOccurred(response["message"].toString());
    }
}

void VpnController::setStatus(const QString &newStatus) {
    if (m_status != newStatus) {
        m_status = newStatus;
        emit statusChanged();
    }
}
