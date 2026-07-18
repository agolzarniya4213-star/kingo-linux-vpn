#pragma once

#include <QObject>
#include <QVariantList>
#include <QString>
#include <QTimer>

class IpcClient;

class AppCore : public QObject {
    Q_OBJECT
    Q_PROPERTY(QString connectionStatus READ connectionStatus NOTIFY statusChanged)
    Q_PROPERTY(QString downloadSpeed READ downloadSpeed NOTIFY statsChanged)
    Q_PROPERTY(QString uploadSpeed READ uploadSpeed NOTIFY statsChanged)
    Q_PROPERTY(QString totalDownload READ totalDownload NOTIFY statsChanged)
    Q_PROPERTY(QString totalUpload READ totalUpload NOTIFY statsChanged)
    Q_PROPERTY(QString selectedServer READ selectedServer WRITE setSelectedServer NOTIFY selectedServerChanged)
    Q_PROPERTY(QVariantList serverList READ serverList NOTIFY serverListChanged)

public:
    explicit AppCore(QObject *parent = nullptr);
    ~AppCore() override;

    QString connectionStatus() const;
    QString downloadSpeed() const;
    QString uploadSpeed() const;
    QString totalDownload() const;
    QString totalUpload() const;
    QString selectedServer() const;
    void setSelectedServer(const QString &server);
    QVariantList serverList() const;

    Q_INVOKABLE void refreshServers();
    Q_INVOKABLE QString formatBytes(double bytes) const;

public slots:
    void connectVPN();
    void disconnectVPN();

signals:
    void statusChanged();
    void statsChanged();
    void selectedServerChanged();
    void serverListChanged();

private:
    IpcClient *m_ipc = nullptr;
    QTimer *m_timer = nullptr;
    QString m_connectionStatus;
    QString m_downloadSpeed;
    QString m_uploadSpeed;
    QString m_totalDownload;
    QString m_totalUpload;
    QString m_selectedServer;
    QVariantList m_serverList;

    void pollStatus();
};
