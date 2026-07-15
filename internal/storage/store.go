package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kingo-linux/vpn/internal/model"
)

type Store struct {
	Servers []model.Server `json:"servers"`
}

func Load(path string) (Store, error) {
	var s Store
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Store{Servers: []model.Server{}}, nil
		}
		return Store{}, err
	}
	if len(data) == 0 {
		return Store{Servers: []model.Server{}}, nil
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return Store{}, err
	}
	if s.Servers == nil {
		s.Servers = []model.Server{}
	}
	return s, nil
}

func Save(path string, s Store) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
