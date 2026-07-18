#pragma once
#include <QObject>
#include <QVariantList>
#include "ipcclient.h"

class VpnController : public QObject {
    Q_OBJECT
    Q_PROPERTY(QString status READ status NOTIFY statusChanged)
    Q_PROPERTY(bool connected READ connected NOTIFY statusChanged)
    Q_PROPERTY(QVariantList servers READ servers NOTIFY serversChanged)

public:
    explicit VpnController(QObject *parent = nullptr);

    QString status() const { return m_status; }
    bool connected() const { return m_status == "connected"; }
    QVariantList servers() const { return m_servers; }

    Q_INVOKABLE void connectToServer(const QString &uri);
    Q_INVOKABLE void disconnectVpn();
    Q_INVOKABLE void refreshStatus();
    Q_INVOKABLE void fetchServers();
    Q_INVOKABLE void addSubscription(const QString &url);
    Q_INVOKABLE void testLatency();

signals:
    void statusChanged();
    void serversChanged();
    void errorOccurred(const QString &error);

private slots:
    void onResponseReceived(const QJsonObject &response);

private:
    IpcClient *m_client;
    QString m_status = "disconnected";
    QVariantList m_servers;
    void setStatus(const QString &newStatus);
};
