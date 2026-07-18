#include <QApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QQuickStyle>
#include "vpncontroller.h"
#include "trayicon.h"

int main(int argc, char *argv[]) {
    // تغییر از QGuiApplication به QApplication برای پشتیبانی از QSystemTrayIcon و منوها
    QApplication app(argc, argv);
    QQuickStyle::setStyle("Basic");

    QQmlApplicationEngine engine;
    
    VpnController controller;
    engine.rootContext()->setContextProperty("vpnController", &controller);

    TrayIcon tray;
    engine.rootContext()->setContextProperty("trayIcon", &tray);
    tray.show();

    engine.load(QUrl(QStringLiteral("qrc:/qt/qml/KingoVPN/main.qml")));
    if (engine.rootObjects().isEmpty()) return -1;

    return app.exec();
}
