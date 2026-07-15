#include <QApplication>
#include "trayicon.h"

int main(int argc, char *argv[]) {
    // در برنامه‌های Tray، برنامه نباید با بسته شدن پنجره متوقف شود
    QApplication app(argc, argv);
    app.setQuitOnLastWindowClosed(false);

    // راه‌اندازی Tray Icon
    TrayIcon tray;

    return app.exec();
}
