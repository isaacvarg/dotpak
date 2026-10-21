# Mise entries run mise use -g

A mise install command is one tool request, such as `node` or `node@20`, passed through unchanged. Install runs a single `mise use -g` as the current user, which installs those versions and records them in the global mise config so a fresh machine can resolve the tools. `mise install` was rejected because an unconfigured tool stays off `PATH`; the cost is that this writes `~/.config/mise/config.toml`, which a dotfiles manager may also own.
