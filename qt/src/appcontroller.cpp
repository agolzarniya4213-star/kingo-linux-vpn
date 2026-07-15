#include "appcontroller.h"

#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonValue>
#include <QLocalSocket>
#include <QStandardPaths>

static QVariantMap parseJsonObject(const QJsonObject &obj) {
    QVariantMap out;
    for (auto it = obj.begin(); it != obj.end(); ++it) {
        if (it.value().isObject()) {
            out.insert(it.key(), parseJsonObject(it.value().toObject()));
        } else if (it.value().isArray()) {
            QVariantList list;
            for (const auto &v : it.value().toArray())
                list.append(v.toVariant());
            out.insert(it.key(), list);
        } else {
            out.insert(it.key(), it.value().toVariant());
        }
    }
    return out;
}

AppController::AppController(QObject *parent) : QObject(parent) {
    const auto configDir = QStandardPaths::writableLocation(QStandardPaths::AppConfigLocation);
    m_socketPath = configDir + "/run/daemon.sock";
    m_state = "unknown";
}

QString AppController::state() const { return m_state; }
QString AppController::socketPath() const { return m_socketPath; }
QString AppController::lastError() const { return m_lastError; }

void AppController::reloadStatus() {
    applyStatus(callDaemon("status"));
}

void AppController::refresh() {
    applyStatus(callDaemon("refresh"));
}

void AppController::stopEngine() {
    applyStatus(callDaemon("stop"));
}

void AppController::startServer(const QVariantMap &server) {
    applyStatus(callDaemon("start", server));
}

void AppController::startFromJson(const QString &jsonText) {
    const auto bytes = jsonText.toUtf8();
    const auto doc = QJsonDocument::fromJson(bytes);
    if (!doc.isObject()) {
        setLastError("Invalid JSON payload");
        return;
    }
    startServer(parseJsonObject(doc.object()));
}

QVariantMap AppController::callDaemon(const QString &action, const QVariantMap &payload) {
    QLocalSocket socket;
    socket.connectToServer(m_socketPath);
    if (!socket.waitForConnected(1500)) {
        return QVariantMap{{"ok", false}, {"error", "daemon not reachable"}, {"state", m_state}};
    }

    QJsonObject req;
    req["action"] = action;
    for (auto it = payload.begin(); it != payload.end(); ++it)
        req[it.key()] = QJsonValue::fromVariant(it.value());

    socket.write(QJsonDocument(req).toJson(QJsonDocument::Compact));
    socket.write("\n");
    socket.flush();

    if (!socket.waitForReadyRead(2000))
        return QVariantMap{{"ok", false}, {"error", "timeout waiting for daemon"}, {"state", m_state}};

    const auto doc = QJsonDocument::fromJson(socket.readAll());
    if (!doc.isObject())
        return QVariantMap{{"ok", false}, {"error", "invalid daemon response"}, {"state", m_state}};

    return parseJsonObject(doc.object());
}

void AppController::applyStatus(const QVariantMap &status) {
    const auto newState = status.value("state").toString();
    if (newState != m_state) {
        m_state = newState;
        emit stateChanged();
    }

    setLastError(status.value("error").toString());
    emit statusChanged(status);
    emit serversChanged(status.value("servers").toList());
}

void AppController::setLastError(const QString &error) {
    if (error == m_lastError)
        return;
    m_lastError = error;
    emit statusChanged(QVariantMap{{"ok", false}, {"error", m_lastError}, {"state", m_state}});
}
