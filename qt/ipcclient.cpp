#include "ipcclient.h"

IpcClient::IpcClient(QObject *parent) : QObject(parent) {}

QString IpcClient::sendCommand(const QString &command) {
    QTcpSocket socket;
    socket.connectToHost("127.0.0.1", 9876);
    if (!socket.waitForConnected(1000)) {
        return "";
    }
    socket.write(command.toUtf8());
    socket.waitForBytesWritten(1000);
    if (socket.waitForReadyRead(1000)) {
        return QString::fromUtf8(socket.readAll());
    }
    return "";
}
