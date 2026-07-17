#include <QApplication>
#include <QStyle>
#include "trayicon.h"

int main(int argc, char *argv[]) {
    // Set Application Metadata for Wayland/Linux Desktop integration
    QCoreApplication::setApplicationName("Kingo VPN");
    QCoreApplication::setApplicationVersion("1.0.0");
    QCoreApplication::setOrganizationName("Kingo");
    
    // Critical fix for Wayland tray icon support
    qputenv("QT_WAYLAND_DESKTOP_FILE_NAME", "kingo-linux-vpn-tray.desktop");

    QApplication app(argc, argv);
    
    // Prevents app from closing when menus are closed
    app.setQuitOnLastWindowClosed(false);

    // Initialize Tray Icon
    TrayIcon tray;

    // Safe Quit connection
    QObject::connect(&tray, &TrayIcon::quitRequested, &app, &QApplication::quit);

    return app.exec();
}
