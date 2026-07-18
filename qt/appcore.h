#ifndef APPCORE_H
#define APPCORE_H

#include <QObject>
#include <QString>
#include <QTimer>

class IpcClient;

class AppCore : public QObject
{
    Q_OBJECT
    Q_PROPERTY(bool isConnected READ isConnected NOTIFY connectionStateChanged)
    Q_PROPERTY(QString statusText READ statusText NOTIFY connectionStateChanged)
    Q_PROPERTY(QString selectedServer READ selectedServer WRITE setSelectedServer NOTIFY selectedServerChanged)
    
    // New Properties for Live Stats
    Q_PROPERTY(QString rxData READ rxData NOTIFY statsUpdated)
    Q_PROPERTY(QString txData READ txData NOTIFY statsUpdated)
    Q_PROPERTY(QString pingData READ pingData NOTIFY statsUpdated)

public:
    explicit AppCore(QObject *parent = nullptr);

    bool isConnected() const;
    QString statusText() const;
    QString selectedServer() const;
    void setSelectedServer(const QString &server);
    
    QString rxData() const;
    QString txData() const;
    QString pingData() const;

    Q_INVOKABLE void toggleConnection();
    Q_INVOKABLE void openServerList();

signals:
    void connectionStateChanged();
    void selectedServerChanged();
    void showServerListRequested();
    void statsUpdated();

private slots:
    void onConnectSuccess(const QString & /*message*/);
    void onDisconnectSuccess(const QString & /*message*/);
    void onCommandError(const QString & /*errorMessage*/);
    void updateMockStats();

private:
    IpcClient *m_ipc;
    bool m_isConnected;
    QString m_statusText;
    QString m_selectedServer;
    
    QTimer *m_statsTimer;
    double m_mockRx;
    double m_mockTx;
    int m_mockPing;
    
    QString formatBytes(double bytes) const;
};
#endif // APPCORE_H
