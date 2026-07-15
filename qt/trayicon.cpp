#include "trayicon.h"
#include <QApplication>
#include <QStyle>
#include <QPainter>
#include <QSvgRenderer>

TrayIcon::TrayIcon(QObject *parent)
    : QObject(parent)
    , m_trayIcon(new QSystemTrayIcon(this))
    , m_menu(new QMenu())
    , m_connectAction(new QAction(tr("اتصال به آخرین سرور"), this))
    , m_disconnectAction(new QAction(tr("قطع اتصال"), this))
    , m_separator(new QAction(this))
    , m_showAction(new QAction(tr("نمایش پنجره اصلی"), this))
    , m_quitAction(new QAction(tr("خروج"), this))
    , m_isConnected(false)
{
    setupIcons();
    setupMenu();
    
    // Connect signals
    connect(m_trayIcon, &QSystemTrayIcon::activated, 
            this, &TrayIcon::handleActivated);
    connect(m_connectAction, &QAction::triggered, 
            this, &TrayIcon::handleConnectAction);
    connect(m_disconnectAction, &QAction::triggered, 
            this, &TrayIcon::handleDisconnectAction);
    connect(m_showAction, &QAction::triggered, 
            this, &TrayIcon::showMainWindowRequested);
    connect(m_quitAction, &QAction::triggered, 
            this, &TrayIcon::quitRequested);
    
    // Initial state
    updateIcon();
    m_trayIcon->show();
    
    // Timer for connecting animation
    m_connectionTimer.setInterval(500);
    connect(&m_connectionTimer, &QTimer::timeout, [this]() {
        static int index = 0;
        QStringList connectingIcons = {
            ":/icons/connecting-symbolic.svg",
            ":/icons/disconnected-symbolic.svg"
        };
        
        if (index % 2 == 0) {
            m_trayIcon->setIcon(QIcon(connectingIcons[0]));
        } else {
            m_trayIcon->setIcon(QIcon(connectingIcons[1]));
        }
        index++;
    });
}

void TrayIcon::setupIcons()
{
    // Set default icon
    m_trayIcon->setIcon(QIcon(":/icons/disconnected-symbolic.svg"));
    m_trayIcon->setToolTip("Kingo VPN - قطع");
    
    // Set context menu
    m_trayIcon->setContextMenu(m_menu);
}

void TrayIcon::setupMenu()
{
    // Add actions to menu
    m_menu->addAction(m_connectAction);
    m_menu->addAction(m_disconnectAction);
    
    m_separator->setSeparator(true);
    m_menu->addAction(m_separator);
    
    m_menu->addAction(m_showAction);
    m_menu->addAction(m_quitAction);
    
    // Initial states
    m_disconnectAction->setEnabled(false);
    m_connectAction->setEnabled(true);
}

void TrayIcon::handleActivated(QSystemTrayIcon::ActivationReason reason)
{
    if (reason == QSystemTrayIcon::Trigger) {
        // Single click: toggle connection or show main window
        if (m_isConnected) {
            showMainWindowRequested();
        } else {
            handleConnectAction();
        }
    }
}

void TrayIcon::handleConnectAction()
{
    if (!m_currentServer.isEmpty()) {
        emit connectRequested(m_currentServer);
    } else {
        showNotification("خطا", "لطفاً ابتدا یک سرور انتخاب کنید", 
                        QSystemTrayIcon::Warning);
        emit showMainWindowRequested();
    }
}

void TrayIcon::handleDisconnectAction()
{
    emit disconnectRequested();
}

void TrayIcon::setConnected(bool connected, const QString &serverName)
{
    m_isConnected = connected;
    m_currentServer = serverName;
    
    // Stop connecting animation if running
    if (m_connectionTimer.isActive()) {
        m_connectionTimer.stop();
    }
    
    // Update UI
    m_connectAction->setEnabled(!connected);
    m_disconnectAction->setEnabled(connected);
    
    if (connected) {
        m_connectAction->setText(tr("اتصال مجدد به %1").arg(serverName));
        showNotification("اتصال برقرار شد", 
                        QString("به سرور %1 متصل شدید").arg(serverName),
                        QSystemTrayIcon::Information);
    } else {
        m_connectAction->setText(tr("اتصال به آخرین سرور"));
    }
    
    updateIcon();
    emit connectionStateChanged(connected);
    emit currentServerChanged(serverName);
}

void TrayIcon::setConnecting(const QString &serverName)
{
    m_currentServer = serverName;
    m_isConnected = false;
    
    // Start connecting animation
    m_connectionTimer.start();
    
    // Update UI
    m_connectAction->setEnabled(false);
    m_disconnectAction->setEnabled(true);
    
    m_trayIcon->setToolTip(QString("در حال اتصال به %1...").arg(serverName));
    emit connectionStateChanged(false);
}

void TrayIcon::setError(const QString &errorMessage)
{
    m_connectionTimer.stop();
    updateIcon();
    
    showNotification("خطا در اتصال", errorMessage, 
                    QSystemTrayIcon::Critical);
}

void TrayIcon::updateServerList(const QStringList &servers)
{
    // Here you could implement dynamic server selection menu
    // For now, we just update the last used server
    if (!servers.isEmpty()) {
        m_currentServer = servers.first();
    }
}

void TrayIcon::updateIcon()
{
    QString iconPath;
    
    if (m_isConnected) {
        iconPath = ":/icons/connected-symbolic.svg";
        m_trayIcon->setToolTip(QString("Kingo VPN - متصل به %1").arg(m_currentServer));
    } else {
        iconPath = ":/icons/disconnected-symbolic.svg";
        m_trayIcon->setToolTip("Kingo VPN - قطع");
    }
    
    m_trayIcon->setIcon(QIcon(iconPath));
}

void TrayIcon::showNotification(const QString &title, const QString &message, 
                               QSystemTrayIcon::MessageIcon icon)
{
    if (m_trayIcon->isVisible()) {
        m_trayIcon->showMessage(title, message, icon, 3000);
    }
}

TrayIcon::~TrayIcon()
{
    m_trayIcon->hide();
}
