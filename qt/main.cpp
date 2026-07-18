#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QQuickStyle>
#include "vpncontroller.h"

int main(int argc, char *argv[]) {
    QGuiApplication app(argc, argv);
    QQuickStyle::setStyle("Basic");

    QQmlApplicationEngine engine;
    
    VpnController controller;
    engine.rootContext()->setContextProperty("vpnController", &controller);

    engine.load(QUrl(QStringLiteral("qrc:/qt/qml/KingoVPN/main.qml")));
    if (engine.rootObjects().isEmpty()) return -1;

    return app.exec();
}
