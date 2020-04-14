# Vault UUID Secrets Plugin

UUID is an example secrets engine plugin for [HashiCorp
Vault](https://www.vaultproject.io/). It is meant for demonstration purposes
only and should never be used in production.

This engine was copied and modified from the example secrets/mock plugin from
the [hashicorp/vault-guides][guides] repository.

## Usage

All commands can be run using the provided [Makefile](./Makefile). However, it
may be instructive to look at the commands to gain a greater understanding of
how Vault registers plugins. Using the Makefile will result in running the Vault
server in `dev` mode. Do not run Vault in `dev` mode in production. The `dev`
server allows you to configure the plugin directory as a flag, and automatically
registers plugin binaries in that directory. In production, plugin binaries must
be manually registered.

This will build the plugin binary and start the Vault dev server:
```
# Build UUID plugin and start Vault dev server with plugin automatically registered
$ make all
```

## TODO: update below docs to be relevant for UUID
Now open a new terminal window and run the following commands:
```
# Open a new terminal window and export Vault dev server http address
$ export VAULT_ADDR='http://127.0.0.1:8200'

# Enable the UUID plugin
$ make enable

# Write a secret to the UUID secrets engine
$ vault write uuid/new
Success! 
<example output>

```

[guides]: https://github.com/hashicorp/vault-guides
