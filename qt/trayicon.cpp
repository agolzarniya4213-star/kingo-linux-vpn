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
    setupUI();

    connect(m_trayIcon, &QSystemTrayIcon::activated, 
            this, &TrayIcon::handleActivated);
            
    connect(m_connectAction, &QAction::triggered, this, [this]() {
        m_ipcClient->sendCommand("connect");
    });
    
    connect(m_disconnectAction, &QAction::triggered, this, [this]() {
        m_ipcClient->sendCommand("disconnect");
    });
    
    connect(m_quitAction, &QAction::triggered, 
            this, &TrayIcon::quitRequested);

    connect(m_ipcClient, &IpcClient::commandSuccess, this, [this](const QString &msg) {
        bool isNowConnected = !m_connectAction->isEnabled();
        if (isNowConnected) {
            onDisconnectSuccess(msg);
        } else {
            onConnectSuccess(msg);
        }
    });
    
    connect(m_ipcClient, &IpcClient::commandError, 
            this, &TrayIcon::onCommandError);

    m_trayIcon->show();
}

TrayIcon::~TrayIcon() {
    m_trayIcon->hide();
}

void TrayIcon::setupUI()
{
    m_trayIcon->setIcon(QApplication::style()->standardIcon(QStyle::SP_ComputerIcon));
    m_trayIcon->setToolTip("Kingo VPN - Disconnected");
    m_trayIcon->setContextMenu(m_menu);

    m_menu->addAction(m_connectAction);
    m_menu->addAction(m_disconnectAction);
    m_menu->addSeparator();
    m_menu->addAction(m_showAction);
    m_menu->addAction(m_quitAction);

    m_disconnectAction->setEnabled(false);
}

void TrayIcon::handleActivated(QSystemTrayIcon::ActivationReason reason)
{
    if (reason == QSystemTrayIcon::Trigger) {
        // Future: Toggle main window
    }
}

void TrayIcon::onConnectSuccess(const QString &message) {
    updateVisualState(true);
    m_trayIcon->showMessage("Kingo VPN", message, QSystemTrayIcon::Information, 3000);
}

void TrayIcon::onDisconnectSuccess(const QString &message) {
    updateVisualState(false);
    m_trayIcon->showMessage("Kingo VPN", message, QSystemTrayIcon::Information, 3000);
}

void TrayIcon::onCommandError(const QString &errorMessage) {
    m_trayIcon->showMessage("Kingo VPN Error", errorMessage, QSystemTrayIcon::Critical, 5000);
}

void TrayIcon::updateVisualState(bool isConnected) {
    if (isConnected) {
        m_trayIcon->setIcon(QApplication::style()->standardIcon(QStyle::SP_DialogApplyButton));
        m_trayIcon->setToolTip("Kingo VPN - Connected");
        m_connectAction->setEnabled(false);
        m_disconnectAction->setEnabled(true);
    } else {
        m_trayIcon->setIcon(QApplication::style()->standardIcon(QStyle::SP_ComputerIcon));
        m_trayIcon->setToolTip("Kingo VPN - Disconnected");
        m_connectAction->setEnabled(true);
        m_disconnectAction->setEnabled(false);
    }
}
