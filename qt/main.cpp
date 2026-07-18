#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QQuickStyle>
#include "appcore.h"
#include "ipcclient.h"
#include <QProcess>
#include <QCoreApplication>
#include <QFileInfo>

int main(int argc, char *argv[])
{
    QGuiApplication app(argc, argv);
    app.setQuitOnLastWindowClosed(false);

    // Apply Material Design Theme (Like Hiddify)
    QQuickStyle::setStyle("Material");

    // Auto-start daemon safely
    QProcess daemonProcess;
    QString daemonPath = QCoreApplication::applicationDirPath() + "/kingo-daemon";
    
    if (!QFileInfo::exists(daemonPath)) {
        qCritical() << "Daemon binary missing at:" << daemonPath;
    } else {
        daemonProcess.start(daemonPath, QStringList());
        if (!daemonProcess.waitForStarted(2000)) {
            qWarning() << "Daemon failed to start.";
        }
    }
    
    QObject::connect(&app, &QGuiApplication::aboutToQuit, [&daemonProcess]() {
        daemonProcess.terminate();
        daemonProcess.waitForFinished(1000);
    });

    // Expose C++ Backend to QML
    AppCore appCore;
    
    QQmlApplicationEngine engine;
    engine.rootContext()->setContextProperty("appCore", &appCore);
    
    const QUrl url(QStringLiteral("qrc:/qml/Main.qml"));
    engine.load(url);
    
    if (engine.rootObjects().isEmpty())
        return -1;

    return app.exec();
}
