#include "vpncontroller.h"
#include <QDateTime>
#include <QJsonObject>
#include <QRandomGenerator>

VpnController::VpnController(QObject *parent)
    : QObject(parent), m_client(new IpcClient(this)) {
    connect(m_client, &IpcClient::responseReceived, this, &VpnController::onResponseReceived);
    connect(m_client, &IpcClient::errorOccurred, this, [this](const QString &err) {
        appendLog("IPC Error: " + err);
        emit errorOccurred(err);
    });

    m_timer = new QTimer(this);
    m_timer->setInterval(1000);
    connect(m_timer, &QTimer::timeout, this, &VpnController::updateConnectionTime);

    refreshStatus();
    fetchServers();
}

QVariantList VpnController::filterServers(const QString &category) const {
    QVariantList filtered;
    for (const auto &srv : m_allServers) {
        if (srv.toMap().value("category").toString() == category) {
            filtered.append(srv);
        }
    }
    return filtered;
}

void VpnController::updateConnectionTime() {
    if (m_status == "connected" && m_connectTime > 0) {
        qint64 elapsed = QDateTime::currentSecsSinceEpoch() - m_connectTime;
        int h = elapsed / 3600;
        int m = (elapsed % 3600) / 60;
        int s = elapsed % 60;
        m_connectionTime = QString::asprintf("%02d:%02d:%02d", h, m, s);
        emit timeChanged();
    }
}

void VpnController::appendLog(const QString &log) {
    m_logs += QDateTime::currentDateTime().toString("hh:mm:ss") + " - " + log + "\n";
    emit logsChanged();
}

QString VpnController::generateRequestID() {
    return QString::number(QDateTime::currentMSecsSinceEpoch()) + "-" +
           QString::number(QRandomGenerator::global()->generate());
}

void VpnController::connectToServer(const QString &uri) {
    appendLog("Connecting to server...");
    setStatus("connecting");
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "connect_server";
    req["server_uri"] = uri;
    m_client->sendRequest(req);
}

void VpnController::autoConnect() {
    appendLog("Finding best server...");
    setStatus("connecting");
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "auto_connect";
    m_client->sendRequest(req);
}

void VpnController::disconnectVpn() {
    appendLog("Disconnecting...");
    setStatus("disconnecting");
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
    appendLog("Updating subscription...");
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "add_subscription";
    req["sub_url"] = url;
    m_client->sendRequest(req);
}

void VpnController::clearServers() {
    appendLog("Clearing servers...");
    m_allServers.clear();
    emit serversChanged();

    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "clear_servers";
    m_client->sendRequest(req);
}

void VpnController::testLatency() {
    appendLog("Testing latency...");
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

void VpnController::getIP() {
    QJsonObject req;
    req["request_id"] = generateRequestID();
    req["action"] = "get_ip";
    m_client->sendRequest(req);
}

void VpnController::onResponseReceived(const QJsonObject &response) {
    if (response.contains("state")) {
        setStatus(response["state"].toString());
    } else if (response.contains("status")) {
        setStatus(response["status"].toString());
    }

    if (response.contains("servers")) {
        m_allServers = response["servers"].toVariant().toList();
        emit serversChanged();
        appendLog("Server list updated (" + QString::number(m_allServers.size()) + " items).");
    }

    if (response.contains("upload") && response.contains("download")) {
        m_uploadSpeed = response["upload"].toVariant().toLongLong();
        m_downloadSpeed = response["download"].toVariant().toLongLong();
        emit trafficChanged();
    }

    if (response.contains("ip") && !response["ip"].toString().isEmpty()) {
        m_ipAddress = response["ip"].toString();
        emit ipChanged();
    }

    if (response.contains("success") && !response["success"].toBool()) {
        QString message = response.value("message").toString();
        if (!message.isEmpty()) {
            appendLog("Error: " + message);
            emit errorOccurred(message);
        }
        if (m_status == "connecting" || m_status == "disconnecting") {
            setStatus("disconnected");
        }
    }
}

void VpnController::setStatus(const QString &newStatus) {
    if (m_status == newStatus) return;

    m_status = newStatus;
    emit statusChanged();
    appendLog("Status changed to: " + newStatus.toUpper());

    if (newStatus == "connected") {
        m_connectTime = QDateTime::currentSecsSinceEpoch();
        m_timer->start();
        getIP();
    } else {
        m_timer->stop();
        m_connectionTime = "00:00:00";
        m_ipAddress = "0.0.0.0";
        emit timeChanged();
        emit ipChanged();
    }
}
