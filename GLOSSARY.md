# Dotpak

Dotpak records the packages, apps, and plugins a machine should have, and installs them by group.

## Language

**Entry**:
A named package, app, or plugin in the manifest. It belongs to one group and one install type.
_Avoid_: package, item

**Group**:
A named set of entries.
_Avoid_: tag, profile

**all**:
The reserved group included in every install. It is not a group the user creates.
_Avoid_: default group

**Install type**:
The backend responsible for an entry: pacman, aur, flatpak, omarchy, or mise.
_Avoid_: source, package manager

**Install command**:
The identifier an entry asks its install type to install. For mise, a bare tool name means latest, and a versioned name such as `node@20` means that version.
_Avoid_: shell command
