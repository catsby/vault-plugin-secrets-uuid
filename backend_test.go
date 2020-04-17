package uuid

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/acctest"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/y0ssar1an/q"
)

var testHelper *acctest.Helper

// // RunningAsPlugin returns true if it detects the usual Terraform plugin
// // detection environment variables, suggesting that the current process is
// // being launched as a plugin server.
// // TODO: move to acctest package
// func RunningAsPlugin() bool {
// 	magicCookieKey :=
// 		"VAULT_BACKEND_PLUGIN"
// 	magicCookieValue :=
// 		"6669da05-b1c8-4f49-97d9-c8e5bed98e20"

// 	rap := os.Getenv(magicCookieKey) == magicCookieValue

// 	return os.Getenv(magicCookieKey) == magicCookieValue
// }

func TestMain(m *testing.M) {
	q.Q("-->> starting TestMain from plugin")
	if err := acctest.Setup("uuid"); err != nil {
		panic(err)
	}
	acctest.Run(m)
}

func TestAccUUID_Docker(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.Log("VAULT_ACC is not set")
		t.SkipNow()
	}

	if testHelper == nil {
		t.Fatal("nil helper")
	}

	client := testHelper.Client

	err := client.Sys().Mount("uuid", &api.MountInput{
		Type: "uuid",
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := client.Logical().Write("uuid/new", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if s == nil {
		t.Fatal("nil uuid returned")
	}
	if s.Data["uuid"] == "" {
		t.Fatalf("empty data/uuid: %#v", s.Data)
	}
	t.Logf("got uuid: %#v", s.Data["uuid"])
}

// TODO: remove this POC
func TestUUID_Basic(t *testing.T) {
	if os.Getenv("VAULT_ACC") != "" {
		t.Log("VAULT_ACC is not set")
		t.SkipNow()
	}

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"uuid": Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	err := client.Sys().Mount("uuid", &api.MountInput{
		Type: "uuid",
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := client.Logical().Write("uuid/new", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if s == nil {
		t.Fatal("nil uuid returned")
	}
	if s.Data["uuid"] == "" {
		t.Fatalf("empty data/uuid: %#v", s.Data)
	}

	// _, err = client.Logical().Write("transit/keys/foobar", map[string]interface{}{
	// 	"type": "ecdsa-p384",
	// })
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// _, err = client.Logical().Write("transit/keys/bar", map[string]interface{}{
	// 	"type": "ed25519",
	// })
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// _, err = client.Logical().Read("transit/keys/foo")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// _, err = client.Logical().Read("transit/keys/foobar")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// _, err = client.Logical().Read("transit/keys/bar")
	// if err != nil {
	// 	t.Fatal(err)
	// }

}
