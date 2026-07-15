#include "serverlistmodel.h"

#include <QVariantMap>

ServerListModel::ServerListModel(QObject *parent) : QAbstractListModel(parent) {}

int ServerListModel::rowCount(const QModelIndex &parent) const {
    return parent.isValid() ? 0 : m_items.size();
}

QVariant ServerListModel::data(const QModelIndex &index, int role) const {
    if (!index.isValid() || index.row() < 0 || index.row() >= m_items.size())
        return {};

    const auto &item = m_items.at(index.row());
    switch (role) {
    case NameRole: return item.name;
    case ProtocolRole: return item.protocol;
    case ConfigRole: return item.config;
    case PingRole: return item.pingMs;
    case FavoriteRole: return item.favorite;
    case GroupRole: return item.group;
    default: return {};
    }
}

QHash<int, QByteArray> ServerListModel::roleNames() const {
    return {
        {NameRole, "name"},
        {ProtocolRole, "protocol"},
        {ConfigRole, "config"},
        {PingRole, "pingMs"},
        {FavoriteRole, "favorite"},
        {GroupRole, "group"},
    };
}

void ServerListModel::setServers(const QVariantList &servers) {
    beginResetModel();
    m_items.clear();

    for (const auto &v : servers) {
        const auto m = v.toMap();
        Item item;
        item.name = m.value("name").toString();
        item.protocol = m.value("protocol").toString();
        item.config = m.value("config").toString();
        item.pingMs = m.value("ping_ms").toInt(-1);
        item.favorite = m.value("favorite").toBool();
        item.group = m.value("group").toString("servers");
        m_items.push_back(item);
    }

    endResetModel();
}

QVariantMap ServerListModel::get(int row) const {
    QVariantMap out;
    if (row < 0 || row >= m_items.size())
        return out;

    const auto &item = m_items.at(row);
    out["name"] = item.name;
    out["protocol"] = item.protocol;
    out["config"] = item.config;
    out["ping_ms"] = item.pingMs;
    out["favorite"] = item.favorite;
    out["group"] = item.group;
    return out;
}
