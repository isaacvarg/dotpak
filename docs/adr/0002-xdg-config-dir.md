# State lives in the XDG config directory

User state is read and written under `$XDG_CONFIG_HOME/dotpak`, or `~/.config/dotpak` when that variable is unset, so a dotfiles manager can track one directory. The tool does not take a path to that state.
