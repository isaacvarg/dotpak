// Package groups manages the package groups as user can assign
package groups

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/isaacvarg/dotpak/internal/storage"
)

var ErrDuplicateGroup = errors.New("group already exists")

type Groups struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Name string `json:"name"`
}

func Load() (*Groups, error) {
	data, err := os.ReadFile(groupsFile())
	if errors.Is(err, os.ErrNotExist) {
		return &Groups{}, nil
	}
	if err != nil {
		return nil, err
	}

	var g Groups
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}

	return &g, nil
}

func (g *Groups) Add(name string) error {
	for _, existing := range g.Groups {
		if strings.EqualFold(existing.Name, name) {
			return ErrDuplicateGroup
		}
	}

	g.Groups = append(g.Groups, Group{Name: name})

	return nil
}

func (g *Groups) Save() error {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(storage.ConfigDir(), 0o755); err != nil {
		return err
	}

	return os.WriteFile(groupsFile(), data, 0o0644)
}

func groupsFile() string {
	return path.Join(storage.ConfigDir(), "groups.json")
}
