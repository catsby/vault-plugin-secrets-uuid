package main

import (
	"os"

	uuid "github.com/catsby/vault-plugin-secrets-uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
	"github.com/y0ssar1an/q"
)

func RunningAsPlugin() bool {
	magicCookieKey :=
		"VAULT_BACKEND_PLUGIN"
	magicCookieValue :=
		"6669da05-b1c8-4f49-97d9-c8e5bed98e20"

	return os.Getenv(magicCookieKey) == magicCookieValue
}

func main() {
	if RunningAsPlugin() {
		q.Q("--> running as plugin from cmd/main")
	} else {
		q.Q("--> NOT running as plugin from cmd/main")
	}
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])
	q.Q("os.Args[1:] from main:", os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: uuid.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		logger := hclog.New(&hclog.LoggerOptions{})

		logger.Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}
