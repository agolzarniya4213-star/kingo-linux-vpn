#include "trayicon.h"
#include <QApplication>
#include <QStyle>

TrayIcon::TrayIcon(QObject *parent) : QObject(parent) {
    m_tray = new QSystemTrayIcon(this);
    // استفاده از یک آیکون استاندارد سیستم (در آیا می‌توان آیکون اختصاصی گذاشت)
    m_tray->setIcon(QApplication::style()->standardIcon(QStyle::SP_ComputerIcon));
    m_tray->setToolTip("Kingo Linux VPN");

    m_menu = new QMenu();
    QAction *showAction = m_menu->addAction("Show Window");
    QAction *connectAction = m_menu->addAction("Connect");
    QAction *disconnectAction = m_menu->addAction("Disconnect");
    m_menu->addSeparator();
    QAction *quitAction = m_menu->addAction("Quit");

    m_tray->setContextMenu(m_menu);

    connect(m_tray, &QSystemTrayIcon::activated, this, &TrayIcon::onActivated);
    connect(showAction, &QAction::triggered, this, &TrayIcon::activateRequested);
    connect(connectAction, &QAction::triggered, this, &TrayIcon::connectRequested);
    connect(disconnectAction, &QAction::triggered, this, &TrayIcon::disconnectRequested);
    connect(quitAction, &QAction::triggered, this, &TrayIcon::quitRequested);
}

void TrayIcon::onActivated(QSystemTrayIcon::ActivationReason reason) {
    if (reason == QSystemTrayIcon::Trigger || reason == QSystemTrayIcon::DoubleClick) {
        emit activateRequested();
    }
}
