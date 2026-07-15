#pragma once

#include <QObject>
#include <QString>
#include <QVariantList>
#include <QVariantMap>

class AppController : public QObject {
    Q_OBJECT
    Q_PROPERTY(QString state READ state NOTIFY stateChanged)
    Q_PROPERTY(QString socketPath READ socketPath CONSTANT)
    Q_PROPERTY(QString lastError READ lastError NOTIFY statusChanged)
public:
    explicit AppController(QObject *parent = nullptr);

    QString state() const;
    QString socketPath() const;
    QString lastError() const;

    Q_INVOKABLE void reloadStatus();
    Q_INVOKABLE void refresh();
    Q_INVOKABLE void stopEngine();
    Q_INVOKABLE void startServer(const QVariantMap &server);
    Q_INVOKABLE void startFromJson(const QString &jsonText);

signals:
    void stateChanged();
    void statusChanged(const QVariantMap &status);
    void serversChanged(const QVariantList &servers);

private:
    QString m_state;
    QString m_socketPath;
    QString m_lastError;

    QVariantMap callDaemon(const QString &action, const QVariantMap &payload = {});
    void applyStatus(const QVariantMap &status);
    void setLastError(const QString &error);
};
