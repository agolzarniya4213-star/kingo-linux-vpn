#pragma once
#include <QObject>
#include <QLocalSocket>
#include <QJsonObject>
#include <QByteArray>

class IpcClient : public QObject {
    Q_OBJECT
public:
    explicit IpcClient(QObject *parent = nullptr);
    void sendRequest(const QJsonObject &request);

signals:
    void responseReceived(const QJsonObject &response);
    void errorOccurred(const QString &error);

private slots:
    void onReadyRead();
    void onConnected();

private:
    QLocalSocket *m_socket;
    QByteArray m_buffer;
    QByteArray m_pendingRequest;
};
