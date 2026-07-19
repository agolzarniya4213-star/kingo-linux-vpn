#include "trayicon.h"
#include <QApplication>
#include <QStyle>

TrayIcon::TrayIcon(QObject *parent) : QObject(parent) {
    m_tray = new QSystemTrayIcon(this);
    m_tray->setIcon(QApplication::style()->standardIcon(QStyle::SP_ComputerIcon));
    m_tray->setToolTip("Kingo Linux VPN");

    m_menu = new QMenu();
    m_menu->addAction("Show Window", this, &TrayIcon::activateRequested);
    m_menu->addAction("Connect", this, &TrayIcon::connectRequested);
    m_menu->addAction("Disconnect", this, &TrayIcon::disconnectRequested);
    m_menu->addSeparator();
    m_menu->addAction("Quit", this, &TrayIcon::quitRequested);

    m_tray->setContextMenu(m_menu);
    connect(m_tray, &QSystemTrayIcon::activated, this, &TrayIcon::onActivated);
}

// پاکسازی حافظه منو هنگام تخریب شیء TrayIcon
TrayIcon::~TrayIcon() {
    delete m_menu;
}

void TrayIcon::onActivated(QSystemTrayIcon::ActivationReason reason) {
    if (reason == QSystemTrayIcon::Trigger || reason == QSystemTrayIcon::DoubleClick) {
        emit activateRequested();
    }
}
