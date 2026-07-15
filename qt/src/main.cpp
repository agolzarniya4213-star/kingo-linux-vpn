#include <QCoreApplication>
#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QUrl>

#include "appcontroller.h"
#include "models/serverlistmodel.h"

int main(int argc, char *argv[]) {
    QGuiApplication app(argc, argv);

    AppController controller;
    ServerListModel serverListModel;

    QObject::connect(&controller, &AppController::serversChanged,
                     &serverListModel, &ServerListModel::setServers);

    QQmlApplicationEngine engine;
    engine.rootContext()->setContextProperty("AppController", &controller);
    engine.rootContext()->setContextProperty("ServerListModel", &serverListModel);

    const QUrl url(QStringLiteral("qrc:/KingoVpn/qml/Main.qml"));
    QObject::connect(&engine, &QQmlApplicationEngine::objectCreated,
                     &app, [url](QObject *obj, const QUrl &objUrl) {
        if (!obj && objUrl == url)
            QCoreApplication::exit(-1);
    }, Qt::QueuedConnection);

    engine.load(url);
    return app.exec();
}
