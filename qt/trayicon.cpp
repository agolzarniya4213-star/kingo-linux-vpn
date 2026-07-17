#include "trayicon.h"
#include "ipcclient.h"
#include <QApplication>
#include <QStyle>

TrayIcon::TrayIcon(QObject *parent)
    : QObject(parent)
    , m_trayIcon(new QSystemTrayIcon(this))
    , m_menu(new QMenu())
    , m_connectAction(new QAction(tr("Connect"), this))
    , m_disconnectAction(new QAction(tr("Disconnect"), this))
    , m_showAction(new QAction(tr("Show Window"), this))
    , m_quitAction(new QAction(tr("Quit"), this))
    , m_ipcClient(new IpcClient(this))
{
    setupIcons();
    setupMenu();

    connect(m_trayIcon, &QSystemTrayIcon::activated, 
            this, &TrayIcon::handleActivated);
            
    // Real IPC triggers
    connect(m_connectAction, &QAction::triggered, [this]() {
        m_ipcClient->sendCommand("connect");
    });
    connect(m_disconnectAction, &QAction::triggered, [this]() {
        m_ipcClient->sendCommand("disconnect");
    });
    
    connect(m_quitAction, &QAction::triggered, 
            this, &TrayIcon::quitRequested);

    // Handle IPC responses
    connect(m_ipcClient, &IpcClient::commandSuccess, this, [this](const QString &msg) {
        if (m_connectAction->isEnabled()) {
            onConnectSuccess(msg);
        } else {
            onDisconnectSuccess(msg);
        }
    });
    connect(m_ipcClient, &IpcClient::commandError, this, &TrayIcon::onCommandError);

    m_trayIcon->show();
}

TrayIcon::~TrayIcon() {
    m_trayIcon->hide();
}

void TrayIcon::setupIcons()
{
    m_trayIcon->setIcon(QApplication::style()->standardIcon(QStyle::SP_ComputerIcon));
    m_trayIcon->setToolTip("Kingo VPN - Disconnected");
    m_trayIcon->setContextMenu(m_menu);
}

void TrayIcon::setupMenu()
{
    m_menu->addAction(m_connectAction);
    m_menu->addAction(m_disconnectAction);
    m_menu->addSeparator();
    m_menu->addAction(m_showAction);
    m_menu->addAction(m_quitAction);

    m_disconnectAction->setEnabled(false);
}

void TrayIcon::handleActivated(QSystemTrayIcon::ActivationReason reason)
{
    // Future: handle left click
}

void TrayIcon::onConnectSuccess(const QString &message) {
    m_trayIcon->setIcon(QApplication::style()->standardIcon(QStyle::SP_DialogApplyButton)); // Changes icon
    m_trayIcon->setToolTip("Kingo VPN - Connected");
    m_connectAction->setEnabled(false);
    m_disconnectAction->setEnabled(true);
    m_trayIcon->showMessage("Kingo VPN", message, QSystemTrayIcon::Information, 3000); // Desktop notification
}

void TrayIcon::onDisconnectSuccess(const QString &message) {
    m_trayIcon->setIcon(QApplication::style()->standardIcon(QStyle::SP_ComputerIcon)); // Reverts icon
    m_trayIcon->setToolTip("Kingo VPN - Disconnected");
    m_connectAction->setEnabled(true);
    m_disconnectAction->setEnabled(false);
    m_trayIcon->showMessage("Kingo VPN", message, QSystemTrayIcon::Information, 3000);
}

void TrayIcon::onCommandError(const QString &message) {
    m_trayIcon->showMessage("Kingo VPN Error", message, QSystemTrayIcon::Critical, 3000);
}
