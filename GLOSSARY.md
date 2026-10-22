# Dotpak

Dotpak records the packages, apps, and plugins a machine should have, and installs them by group.

## Language

**Entry**:
A named package, app, or plugin in the manifest. It belongs to one group and one install type, and its name identifies it.
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
The identifier an entry asks its install type to install. It is the entry's name unless a different identifier was recorded, and for mise a bare tool name means latest while a versioned name such as `node@20` means that version.
_Avoid_: shell command, custom command

**Remove**:
Taking an entry out of the manifest. The package, app, or plugin stays on the machine.
_Avoid_: delete, uninstall

**Edit**:
A change to an entry's name, group, install type, or install command.
_Avoid_: update

**Filter**:
A group, an install type, or both, limiting which entries are loaded. An entry is loaded only when it matches every filter that is set.
_Avoid_: query

**Query**:
A fragment of an entry's name. Among the entries a filter has loaded, the list shows an entry only when its name contains the query, ignoring case.
_Avoid_: filter, search
