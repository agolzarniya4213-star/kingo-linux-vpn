#pragma once
#include <QObject>
#include "ipcclient.h"

class VpnController : public QObject {
    Q_OBJECT
    Q_PROPERTY(QString status READ status NOTIFY statusChanged)
    Q_PROPERTY(bool connected READ connected NOTIFY statusChanged)

public:
    explicit VpnController(QObject *parent = nullptr);

    QString status() const { return m_status; }
    bool connected() const { return m_status == "connected"; }

    Q_INVOKABLE void connectVpn(const QString &configPath);
    Q_INVOKABLE void disconnectVpn();
    Q_INVOKABLE void refreshStatus();

signals:
    void statusChanged();

private slots:
    void onResponseReceived(const QJsonObject &response);

private:
    IpcClient *m_client;
    QString m_status = "disconnected";
    void setStatus(const QString &newStatus);
};
