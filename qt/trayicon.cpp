#include "trayicon.h"
#include <QApplication>
#include <QStyle>

TrayIcon::TrayIcon(QObject *parent) : QObject(parent) {
    m_tray = new QSystemTrayIcon(this);
    
    // Try to load custom icon, fallback to system icon
    QIcon customIcon(":/assets/icon.png");
    if (!customIcon.isNull()) {
        m_tray->setIcon(customIcon);
    } else {
        m_tray->setIcon(QApplication::style()->standardIcon(QStyle::SP_ComputerIcon));
    }
    
    m_tray->setToolTip("Kingo VPN v1.4");

    m_menu = new QMenu();
    // FIX: Removed setParent(this) because QMenu requires a QWidget parent, not QObject.
    m_menu->addAction("Show Window", this, &TrayIcon::activateRequested);
    m_menu->addAction("Connect", this, &TrayIcon::connectRequested);
    m_menu->addAction("Disconnect", this, &TrayIcon::disconnectRequested);
    m_menu->addSeparator();
    m_menu->addAction("Quit", this, &TrayIcon::quitRequested);

    m_tray->setContextMenu(m_menu);
    connect(m_tray, &QSystemTrayIcon::activated, this, &TrayIcon::onActivated);
}

TrayIcon::~TrayIcon() {
    delete m_menu; // Memory is safely managed here
}

void TrayIcon::onActivated(QSystemTrayIcon::ActivationReason reason) {
    if (reason == QSystemTrayIcon::Trigger || reason == QSystemTrayIcon::DoubleClick) {
        emit activateRequested();
    }
}
