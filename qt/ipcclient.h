#pragma once

#include <QObject>
#include <QTcpSocket>

class IpcClient : public QObject {
    Q_OBJECT
public:
    explicit IpcClient(QObject *parent = nullptr);
    QString sendCommand(const QString &command);
};
