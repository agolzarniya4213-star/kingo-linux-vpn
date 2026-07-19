package storage

import (
    "database/sql"
    "fmt"
    "log/slog"

    _ "github.com/mattn/go-sqlite3"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
)

type SQLiteStorage struct {
    db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
    db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
    if err != nil {
        return nil, err
    }

    schema := `
    CREATE TABLE IF NOT EXISTS servers (
        id TEXT PRIMARY KEY,
        name TEXT,
        address TEXT,
        port INTEGER,
        protocol TEXT,
        uri TEXT,
        latency INTEGER DEFAULT 0
    );`
    if _, err := db.Exec(schema); err != nil {
        return nil, fmt.Errorf("schema failed: %w", err)
    }
    return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) SaveServers(servers []model.Server) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec("DELETE FROM servers")
    if err != nil {
        return err
    }

    stmt, err := tx.Prepare("INSERT INTO servers(id, name, address, port, protocol, uri, latency) VALUES(?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, srv := range servers {
        _, err = stmt.Exec(srv.ID, srv.Name, srv.Address, srv.Port, srv.Protocol, srv.URI, srv.Latency)
        if err != nil {
            return err
        }
    }
    return tx.Commit()
}

func (s *SQLiteStorage) GetServers() ([]model.Server, error) {
    rows, err := s.db.Query("SELECT id, name, address, port, protocol, uri, latency FROM servers")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var servers []model.Server
    for rows.Next() {
        var srv model.Server
        if err := rows.Scan(&srv.ID, &srv.Name, &srv.Address, &srv.Port, &srv.Protocol, &srv.URI, &srv.Latency); err != nil {
            return nil, err
        }
        servers = append(servers, srv)
    }
    return servers, rows.Err()
}

// FIX BUG-051: Log error on Close
func (s *SQLiteStorage) Close() {
    if err := s.db.Close(); err != nil {
        slog.Error("Failed to close database safely", "error", err)
    }
}
