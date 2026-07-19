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
    // FIX BUG-023: Add WAL mode, busy timeout, and foreign keys
    db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
    if err != nil {
        return nil, err
    }
    
    // FIX BUG-023: Limit connections for SQLite
    db.SetMaxOpenConns(1)

    // FIX BUG-023: Basic schema migration
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
    
    // Set initial schema version if not exists
    db.Exec("INSERT OR IGNORE INTO schema_version (version) VALUES (1)")

    return &SQLiteStorage{db: db}, nil
}

// FIX BUG-022: Use UPSERT instead of DELETE and INSERT
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
            name=excluded.name, 
            address=excluded.address, 
            port=excluded.port, 
            protocol=excluded.protocol, 
            uri=excluded.uri, 
            latency=excluded.latency
    `)
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

func (s *SQLiteStorage) Close() {
    if err := s.db.Close(); err != nil {
        slog.Error("Failed to close database safely", "error", err)
    }
}
