#include "vpncontroller.h"

VpnController::VpnController(QObject *parent) : QObject(parent), m_client(new IpcClient(this)) {
    connect(m_client, &IpcClient::responseReceived, this, &VpnController::onResponseReceived);
    refreshStatus();
}

void VpnController::connectVpn(const QString &configPath) {
    QJsonObject req;
    req["action"] = "connect";
    req["config_path"] = configPath;
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

void VpnController::onResponseReceived(const QJsonObject &response) {
    if (response.contains("state")) {
        setStatus(response["state"].toString());
    }
}

void VpnController::setStatus(const QString &newStatus) {
    if (m_status != newStatus) {
        m_status = newStatus;
        emit statusChanged();
    }
}
