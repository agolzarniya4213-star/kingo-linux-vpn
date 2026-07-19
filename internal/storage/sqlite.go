package storage

import (
    "database/sql"
    "fmt"
    "log/slog"

    _ "github.com/mattn/go-sqlite3"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
)

type SQLiteStorage struct {
    db    *sql.DB
    crypto *CryptoLayer
}

func NewSQLiteStorage(dbPath string, crypto *CryptoLayer) (*SQLiteStorage, error) {
    db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(1)

    schema := `
    CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY);
    CREATE TABLE IF NOT EXISTS servers (
        id TEXT PRIMARY KEY,
        name TEXT,
        address TEXT,
        port INTEGER,
        protocol TEXT,
        uri TEXT,
        latency INTEGER DEFAULT 0,
        last_seen INTEGER DEFAULT 0
    );`
    if _, err := db.Exec(schema); err != nil {
        return nil, fmt.Errorf("schema failed: %w", err)
    }
    db.Exec("INSERT OR IGNORE INTO schema_version (version) VALUES (1)")

    return &SQLiteStorage{db: db, crypto: crypto}, nil
}

func (s *SQLiteStorage) SaveServers(servers []model.Server) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    stmt, err := tx.Prepare(`
        INSERT INTO servers(id, name, address, port, protocol, uri, latency) 
        VALUES(?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET 
            name=excluded.name, address=excluded.address, port=excluded.port, 
            protocol=excluded.protocol, uri=excluded.uri, latency=excluded.latency
    `)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, srv := range servers {
        encryptedURI := s.crypto.Encrypt(srv.URI)
        _, err = stmt.Exec(srv.ID, srv.Name, srv.Address, srv.Port, srv.Protocol, encryptedURI, srv.Latency)
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
        var encryptedURI string
        if err := rows.Scan(&srv.ID, &srv.Name, &srv.Address, &srv.Port, &srv.Protocol, &encryptedURI, &srv.Latency); err != nil {
            return nil, err
        }
        // Decrypt URI
        decryptedURI, err := s.crypto.Decrypt(encryptedURI)
        if err != nil {
            slog.Warn("Failed to decrypt URI, using raw", "error", err)
            decryptedURI = encryptedURI
        }
        srv.URI = decryptedURI
        servers = append(servers, srv)
    }
    return servers, rows.Err()
}

func (s *SQLiteStorage) Close() {
    if err := s.db.Close(); err != nil {
        slog.Error("Failed to close database safely", "error", err)
    }
}
