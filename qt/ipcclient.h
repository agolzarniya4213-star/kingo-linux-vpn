#ifndef IPCCLIENT_H
#define IPCCLIENT_H

#include <QObject>
#include <QLocalSocket>
#include <QJsonObject>

class IpcClient : public QObject
{
    Q_OBJECT
public:
    explicit IpcClient(QObject *parent = nullptr);

    void sendCommand(const QString &action);

signals:
    void commandSuccess(const QString &message);
    void commandError(const QString &errorMessage);

private slots:
    void onSocketError(QLocalSocket::LocalSocketError error);

private:
    QLocalSocket *m_socket;
    QString m_socketPath;

    void initializeSocketPath();
    QJsonObject parseResponse(const QByteArray &data);
};

#endif // IPCCLIENT_H
