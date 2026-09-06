// Package manifest handles storing packages, plugins, and custom installs
package manifest

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/isaacvarg/dotpak/internal/storage"
)


type InstallType string

const (
	InstallPacman InstallType = "pacman" 
	InstallAUR InstallType = "aur"
	InstallFlatpak InstallType = "flatpak"
	InstallMise InstallType = "mise"
	InstallOmarchy InstallType = "omarchy"
)


var ErrDuplicateEntry = errors.New("entry already exists in the manifest")
var ErrNotFound = errors.New("entry not found in the manifest")

func (i InstallType) IsValid() bool {
	switch i {
	case InstallPacman, InstallAUR, InstallFlatpak, InstallMise, InstallOmarchy:
		return true
	}
	return false
}

type Manifest struct {
	Manifest []ManifestEntry `json:"manifest"`
}

type ManifestEntry struct {
	Name string `json:"name"`
	Group string `json:"group"`
	InstallType InstallType `json:"install_type"`
	InstallCommand string `json:"install_command"`
}


func Load() (*Manifest, error) {
	data, err := os.ReadFile(manifestFile())
	if errors.Is(err, os.ErrNotExist) {
		return &Manifest{}, nil
	}
	if err != nil {
		return nil, err
	}

	var m Manifest 
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *Manifest) Save() error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(storage.ConfigDir(), 0o755); err != nil {
		return err
	}

	return os.WriteFile(manifestFile(), data, 0o0644)
}

func (m *Manifest) Add(name string, installType InstallType, group string, installCommand string) error {

	for _, existing := range m.Manifest {
		if strings.EqualFold(existing.Name, name) {
			return ErrDuplicateEntry 
		}
	}

	m.Manifest = append(m.Manifest, ManifestEntry{
		Name: name,
		InstallType: installType,
		Group: group,
		InstallCommand: installCommand,
	})

	return nil
}

func (m *Manifest) Remove(name string) error {
	for i, existing := range m.Manifest {
		if strings.EqualFold(existing.Name, name) {
			m.Manifest = append(m.Manifest[:i], m.Manifest[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}


func manifestFile() string {
	return path.Join(storage.ConfigDir(), "manifest.json")
}
