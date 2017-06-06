
## ircb plugins

### Installing a plugin

  1. Add line to ircb makefile:
	```CGO_ENABLED=1 go build -o skeleton.so -buildmode=plugin github.com/aerth/ircb-plugins/skeleton```

  2. Rebuild ircb if plugin uses newer (or different) version of ircb library
  3. Master command: $plugin myplugin.so


### Gotchas

  * all plugins MUST have unique filenames
  * all plugins MUST have unique package import paths

### Submitting a new plugin

  * Copy 'skeleton' and modify
  * Create a PR


