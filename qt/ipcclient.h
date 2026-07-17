#ifndef IPCCLIENT_H
#define IPCCLIENT_H

#include <QObject>
#include <QLocalSocket>

class IpcClient : public QObject
{
    Q_OBJECT
public:
    explicit IpcClient(QObject *parent = nullptr);
    
    // Async command sender (Does NOT freeze the UI)
    void sendCommand(const QString &action);

signals:
    void commandSuccess(const QString &message);
    void commandError(const QString &message);

private slots:
    void onConnected();
    void onReadyRead();
    void onError(QLocalSocket::LocalSocketError socketError);

private:
    QLocalSocket *m_socket;
    QString m_socketPath;
    QString m_pendingAction;
};

#endif // IPCCLIENT_H
