#include <QApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QQuickStyle>
#include "vpncontroller.h"
#include "trayicon.h"

int main(int argc, char *argv[]) {
    QApplication app(argc, argv);
    QQuickStyle::setStyle("Basic");

    QQmlApplicationEngine engine;
    
    // FIX: Create objects on the heap to prevent premature destruction and QML null references
    VpnController *controller = new VpnController();
    engine.rootContext()->setContextProperty("vpnController", controller);

    TrayIcon *tray = new TrayIcon();
    engine.rootContext()->setContextProperty("trayIcon", tray);
    tray->show();

    engine.load(QUrl(QStringLiteral("qrc:/main.qml")));
    if (engine.rootObjects().isEmpty()) return -1;

    return app.exec();
}
