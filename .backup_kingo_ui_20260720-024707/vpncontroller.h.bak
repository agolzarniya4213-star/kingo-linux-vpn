#pragma once
#include <QObject>
#include <QVariantList>
#include <QGuiApplication>
#include <QClipboard>
#include <QTimer>
#include <QJsonObject>
#include "ipcclient.h"

class VpnController : public QObject {
    Q_OBJECT
    Q_PROPERTY(QString status READ status NOTIFY statusChanged)
    Q_PROPERTY(bool connected READ connected NOTIFY statusChanged)
    Q_PROPERTY(QVariantList servers READ servers NOTIFY serversChanged)
    Q_PROPERTY(QVariantList favorites READ favorites NOTIFY serversChanged)
    Q_PROPERTY(QVariantList custom READ custom NOTIFY serversChanged)
    Q_PROPERTY(qint64 uploadSpeed READ uploadSpeed NOTIFY trafficChanged)
    Q_PROPERTY(qint64 downloadSpeed READ downloadSpeed NOTIFY trafficChanged)
    Q_PROPERTY(QString logs READ logs NOTIFY logsChanged)
    Q_PROPERTY(QString ipAddress READ ipAddress NOTIFY ipChanged)
    Q_PROPERTY(QString connectionTime READ connectionTime NOTIFY timeChanged)

public:
    explicit VpnController(QObject *parent = nullptr);

    QString status() const { return m_status; }
    bool connected() const { return m_status == "connected"; }
    QVariantList servers() const { return filterServers("servers"); }
    QVariantList favorites() const { return filterServers("favorites"); }
    QVariantList custom() const { return filterServers("custom"); }
    qint64 uploadSpeed() const { return m_uploadSpeed; }
    qint64 downloadSpeed() const { return m_downloadSpeed; }
    QString logs() const { return m_logs; }
    QString ipAddress() const { return m_ipAddress; }
    QString connectionTime() const { return m_connectionTime; }

    Q_INVOKABLE void connectToServer(const QString &uri);
    Q_INVOKABLE void autoConnect();
    Q_INVOKABLE void disconnectVpn();
    Q_INVOKABLE void refreshStatus();
    Q_INVOKABLE void fetchServers();
    Q_INVOKABLE void addSubscription(const QString &url);
    Q_INVOKABLE void clearServers();
    Q_INVOKABLE void testLatency();
    Q_INVOKABLE void getTraffic();
    Q_INVOKABLE void getIP();
    Q_INVOKABLE void copyLogs() { QGuiApplication::clipboard()->setText(m_logs); }

signals:
    void statusChanged();
    void serversChanged();
    void trafficChanged();
    void logsChanged();
    void ipChanged();
    void timeChanged();
    void errorOccurred(const QString &error);

private slots:
    void onResponseReceived(const QJsonObject &response);

private:
    IpcClient *m_client;
    QString m_status = "disconnected";
    QVariantList m_allServers;
    qint64 m_uploadSpeed = 0;
    qint64 m_downloadSpeed = 0;
    QString m_logs = "Kingo VPN v0.5 Initialized.\n";
    QString m_ipAddress = "0.0.0.0";
    QString m_connectionTime = "00:00:00";
    qint64 m_connectTime = 0;
    QTimer *m_timer;

    void setStatus(const QString &newStatus);
    void appendLog(const QString &log);
    QString generateRequestID();
    QVariantList filterServers(const QString &category) const;
    void updateConnectionTime();
};
