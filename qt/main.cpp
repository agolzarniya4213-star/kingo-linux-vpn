#include <QApplication>
#include "trayicon.h"

int main(int argc, char *argv[]) {
    // Prevents app from closing when hidden
    QApplication app(argc, argv);
    app.setQuitOnLastWindowClosed(false);

    // Initialize Tray Icon
    TrayIcon tray;

    // Make the "Quit" button actually close the application
    QObject::connect(&tray, &TrayIcon::quitRequested, &app, &QApplication::quit);

    return app.exec();
}
