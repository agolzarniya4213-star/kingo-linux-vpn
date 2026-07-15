#ifndef TRAYICON_H
#define TRAYICON_H

#include <QSystemTrayIcon>
#include <QMenu>
#include <QAction>
#include <QTimer>

class TrayIcon : public QObject
{
    QObject
    Q_PROPERTY(bool isConnected READ isConnected NOTIFY connectionStateChanged)
    Q_PROPERTY(QString currentServer READ currentServer NOTIFY currentServerChanged)

public:
    explicit TrayIcon(QObject *parent = nullptr);
    ~TrayIcon() override;

    // Setters for daemon integration
    void setConnected(bool connected, const QString &serverName = "");
    void setConnecting(const QString &serverName = "");
    void setError(const QString &errorMessage);
    void updateServerList(const QStringList &servers);

signals:
    void connectRequested(const QString &serverName);
    void disconnectRequested();
    void quitRequested();
    void showMainWindowRequested();
    void connectionStateChanged(bool connected);
    void currentServerChanged(const QString &serverName);

private slots:
    void handleActivated(QSystemTrayIcon::ActivationReason reason);
    void handleConnectAction();
    void handleDisconnectAction();

private:
    QSystemTrayIcon *m_trayIcon;
    QMenu *m_menu;
    QAction *m_connectAction;
    QAction *m_disconnectAction;
    QAction *m_separator;
    QAction *m_showAction;
    QAction *m_quitAction;
    
    bool m_isConnected;
    QString m_currentServer;
    QTimer m_connectionTimer;
    
    void setupIcons();
    void setupMenu();
    void updateIcon();
    void showNotification(const QString &title, const QString &message, 
                         QSystemTrayIcon::MessageIcon icon = QSystemTrayIcon::Information);
};

#endif // TRAYICON_H
