#pragma once
#include <QObject>
#include <QSystemTrayIcon>
#include <QMenu>

class TrayIcon : public QObject {
    Q_OBJECT
    Q_PROPERTY(bool visible READ visible WRITE setVisible NOTIFY visibleChanged)
public:
    explicit TrayIcon(QObject *parent = nullptr);
    ~TrayIcon(); // اضافه شدن دستریکتور برای پاکسازی حافظه منو

    bool visible() const { return m_tray->isVisible(); }
    void setVisible(bool v) { m_tray->setVisible(v); emit visibleChanged(); }

    Q_INVOKABLE void show() { m_tray->show(); }
    Q_INVOKABLE void hide() { m_tray->hide(); }
    Q_INVOKABLE void showMessage(const QString &title, const QString &msg) {
        m_tray->showMessage(title, msg);
    }

signals:
    void visibleChanged();
    void activateRequested();
    void connectRequested();
    void disconnectRequested();
    void quitRequested();

private slots:
    void onActivated(QSystemTrayIcon::ActivationReason reason);

private:
    QSystemTrayIcon *m_tray;
    QMenu *m_menu;
};
