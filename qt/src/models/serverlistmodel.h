#pragma once

#include <QAbstractListModel>
#include <QVariantList>

class ServerListModel : public QAbstractListModel {
    Q_OBJECT
public:
    enum Roles {
        NameRole = Qt::UserRole + 1,
        ProtocolRole,
        ConfigRole,
        PingRole,
        FavoriteRole,
        GroupRole
    };

    explicit ServerListModel(QObject *parent = nullptr);

    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    QHash<int, QByteArray> roleNames() const override;

    Q_INVOKABLE void setServers(const QVariantList &servers);
    Q_INVOKABLE QVariantMap get(int row) const;

private:
    struct Item {
        QString name;
        QString protocol;
        QString config;
        int pingMs = -1;
        bool favorite = false;
        QString group = "servers";
    };
    QList<Item> m_items;
};
