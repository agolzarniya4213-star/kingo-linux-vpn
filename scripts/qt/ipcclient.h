#ifndef IPCCLIENT_H
#define IPCCLIENT_H

#include <QObject>
#include <QLocalSocket>

class IpcClient : public QObject
{
    Q_OBJECT
public:
    explicit IpcClient(QObject *parent = nullptr);

    void sendCommand(const QString &action);

signals:
    void commandSuccess(const QString &message);
    void commandError(const QString &message);

private:
    QLocalSocket *m_socket;
    QString m_socketPath;
};

#endif // IPCCLIENT_H
