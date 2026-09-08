# Dotpak

A command-line tool that helps manage packages, apps and plugins for your dotfiles.

`dotpak` lets you declare the packages, apps, and plugins you want installed across machines in a single manifest, organize them into groups (like `work`, `personal`), and install them all with one command.

Dotpak pairs well with a dotfiles manager like [chezmoi](https://www.chezmoi.io/). Add the `~/.config/dotpak` to your chezmoi, create a chezmoi script and using template variables and run `dotpak install -g <group>` to have a new computer ready to go. A more thorough tutorial will come.

Currently works with:
- pacman,
- AUR
- Flatpak
- Omarchy Quattro plugins

## Roadmap

- interactive mode built with bubbletea
- ~~shorthand flags for install type. E.g., -p for pacman~~
- hooks to suggest adding package to dotpak after installation

## Installation

### Requirements

- Arch Linux or an Arch-based distro for pacman/AUR support
- An AUR helper (yay, paru, or pikaur) for aur entries
- Flatpak for flatpak entries
- Omarchy v4.0.0+ for omarchy entries

### AUR
```bash
yay -S dotpak-bin
```

### Release

Download the appropriate archive for your platform from the releases page and place the dotpak binary on your PATH.

### From source

```bash
go install github.com/isaacvarg/dotpak@latest
```

## Usage

*Add an entry to the manifest*

```bash
dotpak add <name> -i <installType> [-g <group>] [-c <installCommand>]
```

- `-i`, `--installType` (required): pacman, aur, flatpak, omarchy, mise
- `-g`, `--group`: group to add this entry to (defaults to all)
- `-c`, `--installCommand`:  the actual command/package id to install, if it
  differs from <name> (e.g. a Flathub app ID or an Omarchy plugin repo)

Shortcuts for the install type can also be used:
- `-p` for pacman
- `-a` for AUR
- `-f` for flatpak
- `-m` for mise
- `-o` for omarchy quattro plugins

Examples:
```bash
dotpak add neovim -i pacman
dotpak add spotify -i flatpak -c com.spotify.Client
dotpak add my-omarchy-theme -i omarchy -c github.com/user/repo -g personal
```

*Install everything in a group*

```bash
dotpak install [-g <group>]
```
Entries in the all group are always included plus the entries in the specified group. This will install each application via the install type

*List manifest entries*

```bash
dotpak list [-g <group>] [-i <installType>]
```

*Remove an entry*
```bash
dotpak remove <name>
```

*Manage groups*
```bash
dotpak group create <name>
dotpak group list
```



