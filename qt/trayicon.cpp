#include "trayicon.h"
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
{
    setupIcons();
    setupMenu();

    connect(m_trayIcon, &QSystemTrayIcon::activated, 
            this, &TrayIcon::handleActivated);
    connect(m_connectAction, &QAction::triggered, 
            this, &TrayIcon::connectRequested);
    connect(m_disconnectAction, &QAction::triggered, 
            this, &TrayIcon::disconnectRequested);
    connect(m_quitAction, &QAction::triggered, 
            this, &TrayIcon::quitRequested);

    m_trayIcon->show();
}

TrayIcon::~TrayIcon() {
    m_trayIcon->hide();
}

void TrayIcon::setupIcons()
{
    // Using standard system icon to prevent file missing errors
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
    if (reason == QSystemTrayIcon::Trigger) {
        // Logic for single click goes here in the next step
    }
}
