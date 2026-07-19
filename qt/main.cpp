#include <QApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QQuickStyle>
#include <QIcon>
#include "vpncontroller.h"
#include "trayicon.h"

int main(int argc, char *argv[]) {
    QApplication app(argc, argv);
    QQuickStyle::setStyle("Basic");
    
    // Set Application Icon
    QIcon appIcon(":/assets/icon.png");
    app.setWindowIcon(appIcon);

    QQmlApplicationEngine engine;
    
    VpnController *controller = new VpnController();
    engine.rootContext()->setContextProperty("vpnController", controller);

    TrayIcon *tray = new TrayIcon(); // Icon is set internally now
    engine.rootContext()->setContextProperty("trayIcon", tray);
    tray->show();

    engine.load(QUrl(QStringLiteral("qrc:/main.qml")));
    if (engine.rootObjects().isEmpty()) return -1;

    return app.exec();
}
