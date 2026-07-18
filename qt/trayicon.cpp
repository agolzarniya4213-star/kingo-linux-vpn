#include "trayicon.h"
#include "appcore.h"
#include <QApplication>
#include <QStyle>
#include <QAction>

TrayIcon::TrayIcon(AppCore *core, QObject *parent)
    : QObject(parent)
    , m_tray(new QSystemTrayIcon(this))
    , m_menu(new QMenu())
    , m_core(core)
{
    m_tray->setIcon(QApplication::style()->standardIcon(QStyle::SP_ComputerIcon));
    m_tray->setToolTip("Kingo VPN");

    auto *connectAction = m_menu->addAction("Connect");
    auto *disconnectAction = m_menu->addAction("Disconnect");
    m_menu->addSeparator();
    auto *quitAction = m_menu->addAction("Quit");

    QObject::connect(connectAction, &QAction::triggered, m_core, &AppCore::connectVPN);
    QObject::connect(disconnectAction, &QAction::triggered, m_core, &AppCore::disconnectVPN);
    QObject::connect(quitAction, &QAction::triggered, qApp, &QApplication::quit);

    m_tray->setContextMenu(m_menu);
    QObject::connect(m_tray, &QSystemTrayIcon::activated, this, &TrayIcon::onTrayActivated);
}

void TrayIcon::show() {
    m_tray->show();
}

void TrayIcon::onTrayActivated(QSystemTrayIcon::ActivationReason reason) {
    if (reason == QSystemTrayIcon::Trigger) {
        // TODO: Toggle main window visibility
    }
}
