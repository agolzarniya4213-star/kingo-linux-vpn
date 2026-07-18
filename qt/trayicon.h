#pragma once

#include <QObject>
#include <QSystemTrayIcon>
#include <QMenu>

class AppCore;

class TrayIcon : public QObject {
    Q_OBJECT
public:
    explicit TrayIcon(AppCore *core, QObject *parent = nullptr);
    void show();

private slots:
    void onTrayActivated(QSystemTrayIcon::ActivationReason reason);

private:
    QSystemTrayIcon *m_tray = nullptr;
    QMenu *m_menu = nullptr;
    AppCore *m_core = nullptr;
};
