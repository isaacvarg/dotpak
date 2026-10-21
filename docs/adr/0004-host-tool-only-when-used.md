# A host tool is required only when the install set uses it

Install runs pacman, then an AUR helper, then Flatpak, then Omarchy, then mise. Each step runs only when the install set contains that install type, so a set with no mise entries does not require mise. Pacman, AUR, and Flatpak are invoked without confirmation prompts, the run stops at the first failure, and that failure exits non-zero.
