#ifndef TRAYICON_H
#define TRAYICON_H

#include <QSystemTrayIcon>
#include <QMenu>

class IpcClient;

class TrayIcon : public QObject
{
    Q_OBJECT
public:
    explicit TrayIcon(QObject *parent = nullptr);
    ~TrayIcon();

signals:
    void quitRequested();

private slots:
    void handleActivated(QSystemTrayIcon::ActivationReason reason);
    void onConnectSuccess(const QString &message);
    void onDisconnectSuccess(const QString &message);
    void onCommandError(const QString &message);

private:
    QSystemTrayIcon *m_trayIcon;
    QMenu *m_menu;
    QAction *m_connectAction;
    QAction *m_disconnectAction;
    QAction *m_showAction;
    QAction *m_quitAction;
    
    IpcClient *m_ipcClient;

    void setupIcons();
    void setupMenu();
};

#endif // TRAYICON_H
