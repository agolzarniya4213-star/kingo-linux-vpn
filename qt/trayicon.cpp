#include "trayicon.h"
#include "ipcclient.h"
#include <QApplication>
#include <QStyle>
#include <QMessageBox>

TrayIcon::TrayIcon(QObject *parent)
    : QObject(parent)
    , m_trayIcon(new QSystemTrayIcon(this))
    , m_menu(new QMenu())
    , m_connectAction(new QAction(tr("Connect"), this))
    , m_disconnectAction(new QAction(tr("Disconnect"), this))
    , m_separator(new QAction(this))
    , m_aboutAction(new QAction(tr("About Kingo VPN"), this))
    , m_quitAction(new QAction(tr("Quit"), this))
    , m_ipcClient(new IpcClient(this))
{
    setupUI();

    // Trigger Async IPC Commands
    connect(m_connectAction, &QAction::triggered, [this]() {
        m_connectAction->setEnabled(false); // Prevent double-clicks
        m_ipcClient->sendCommand("connect");
    });
    
    connect(m_disconnectAction, &QAction::triggered, [this]() {
        m_disconnectAction->setEnabled(false); // Prevent double-clicks
        m_ipcClient->sendCommand("disconnect");
    });

    connect(m_aboutAction, &QAction::triggered, this, &TrayIcon::showAboutDialog);
    connect(m_quitAction, &QAction::triggered, this, &TrayIcon::quitRequested);

    // Handle Async IPC Responses safely
    connect(m_ipcClient, &IpcClient::commandSuccess, this, [this](const QString &msg) {
        bool isNowConnected = !m_connectAction->isEnabled() && m_disconnectAction->isEnabled();
        if (!isNowConnected) {
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

void TrayIcon::setupUI()
{
    m_separator->setSeparator(true);
    
    m_menu->addAction(m_connectAction);
    m_menu->addAction(m_disconnectAction);
    m_menu->addAction(m_separator);
    m_menu->addAction(m_aboutAction);
    m_menu->addAction(m_quitAction);

    m_disconnectAction->setEnabled(false);

    m_trayIcon->setIcon(QApplication::style()->standardIcon(QStyle::SP_ComputerIcon));
    m_trayIcon->setToolTip("Kingo VPN - Disconnected");
    m_trayIcon->setContextMenu(m_menu);
}

void TrayIcon::updateTrayState(bool isConnected)
{
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

void TrayIcon::onConnectSuccess(const QString &message) {
    updateTrayState(true);
    m_trayIcon->showMessage("Kingo VPN", message, QSystemTrayIcon::Information, 3000);
}

void TrayIcon::onDisconnectSuccess(const QString &message) {
    updateTrayState(false);
    m_trayIcon->showMessage("Kingo VPN", message, QSystemTrayIcon::Information, 3000);
}

void TrayIcon::onCommandError(const QString &message) {
    updateTrayState(false); // Reset UI state on error
    m_trayIcon->showMessage("Connection Error", message, QSystemTrayIcon::Critical, 5000);
}

void TrayIcon::showAboutDialog() {
    QMessageBox::information(nullptr, "About Kingo VPN", 
                             "Kingo Linux VPN Client\nVersion: 1.0.0\nArchitecture: Native Qt/Go IPC");
}
