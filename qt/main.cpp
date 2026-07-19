#include <QApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QQuickStyle>
#include <QIcon>
#include "vpncontroller.h"
#include "trayicon.h"

int main(int argc, char *argv[]) {
    QApplication app(argc, argv);
    // FIX: Force Basic style to prevent native desktop widgets from ruining the Material design
    QQuickStyle::setStyle("Basic");
    
    QIcon appIcon(":/assets/icon.png");
    app.setWindowIcon(appIcon);

    QQmlApplicationEngine engine;
    
    VpnController *controller = new VpnController();
    engine.rootContext()->setContextProperty("vpnController", controller);

    TrayIcon *tray = new TrayIcon();
    engine.rootContext()->setContextProperty("trayIcon", tray);
    tray->show();

    engine.load(QUrl(QStringLiteral("qrc:/main.qml")));
    if (engine.rootObjects().isEmpty()) return -1;

    return app.exec();
}
