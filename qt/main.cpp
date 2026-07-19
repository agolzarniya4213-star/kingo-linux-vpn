#include <QApplication> // FIX: Changed from QGuiApplication to QApplication
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QQuickStyle>
#include <QCoreApplication>
#include <QFile>
#include "vpncontroller.h"
#include "trayicon.h"

int main(int argc, char *argv[]) {
    QApplication app(argc, argv);
    QQuickStyle::setStyle("Basic");

    QQmlApplicationEngine engine;
    
    VpnController controller;
    engine.rootContext()->setContextProperty("vpnController", &controller);

    TrayIcon tray;
    engine.rootContext()->setContextProperty("trayIcon", &tray);
    tray.show();

    // 1. Try loading from compiled QRC resources
    engine.load(QUrl(QStringLiteral("qrc:/qt/qml/KingoVPN/main.qml")));
    
    // 2. Fallback: If QRC failed (e.g. dev build issues), load from local file
    if (engine.rootObjects().isEmpty()) {
        QString qmlPath = QCoreApplication::applicationDirPath() + "/../../qt/main.qml";
        if (QFile::exists(qmlPath)) {
            engine.load(QUrl::fromLocalFile(qmlPath));
        }
    }
    
    if (engine.rootObjects().isEmpty()) return -1;

    return app.exec();
}
