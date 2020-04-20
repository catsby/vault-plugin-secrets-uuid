package uuid

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/acctest"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestMain(m *testing.M) {
	// Setup will create the docker cluster, then compile the plugin and register
	// it. Setup will panic if something fails (for now at least)
	acctest.Setup("uuid")

	// acctest.Run wraps the normal m.Run() all with optional call to
	// acctest.TestHelper.Cleanup() to tear down the Docker cluster
	acctest.Run(m)
}

func TestAccUUID_Docker(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.Log("VAULT_ACC is not set")
		t.SkipNow()
	}

	if acctest.TestHelper == nil {
		t.Fatal("nil helper")
	}

	// TODO: function or method to make this safe
	if acctest.TestHelper == nil {
		t.Fatal("expected test helper")
	}
	client := acctest.TestHelper.Client

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

// TestUUID_Basic is an example test using NewTestCluster
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
}
