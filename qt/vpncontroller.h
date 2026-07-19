#pragma once
#include <QObject>
#include <QVariantList>
#include "ipcclient.h"

class VpnController : public QObject {
    Q_OBJECT
    Q_PROPERTY(QString status READ status NOTIFY statusChanged)
    Q_PROPERTY(bool connected READ connected NOTIFY statusChanged)
    Q_PROPERTY(QVariantList servers READ servers NOTIFY serversChanged)
    Q_PROPERTY(qint64 uploadSpeed READ uploadSpeed NOTIFY trafficChanged)
    Q_PROPERTY(qint64 downloadSpeed READ downloadSpeed NOTIFY trafficChanged)

public:
    explicit VpnController(QObject *parent = nullptr);

    QString status() const { return m_status; }
    bool connected() const { return m_status == "connected"; }
    QVariantList servers() const { return m_servers; }
    qint64 uploadSpeed() const { return m_uploadSpeed; }
    qint64 downloadSpeed() const { return m_downloadSpeed; }

    Q_INVOKABLE void connectToServer(const QString &uri);
    Q_INVOKABLE void autoConnect();
    Q_INVOKABLE void disconnectVpn();
    Q_INVOKABLE void refreshStatus();
    Q_INVOKABLE void fetchServers();
    Q_INVOKABLE void addSubscription(const QString &url);
    Q_INVOKABLE void testLatency();
    Q_INVOKABLE void getTraffic();

signals:
    void statusChanged();
    void serversChanged();
    void trafficChanged();
    void errorOccurred(const QString &error);

private slots:
    void onResponseReceived(const QJsonObject &response);

private:
    IpcClient *m_client;
    QString m_status = "disconnected";
    QVariantList m_servers;
    qint64 m_uploadSpeed = 0;
    qint64 m_downloadSpeed = 0;
    void setStatus(const QString &newStatus);
};
