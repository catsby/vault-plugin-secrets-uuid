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

func TestMain(m *testing.M) {
	q.Q("TestMain start")
	stat := m.Run()
	q.Q("TestMain end")
	os.Exit(stat)
}

// TODO: remove this POC
func TestUUID_Basic(t *testing.T) {
	q.Q("--> starting normal test")
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
	q.Q("<-- end normal test")
}

func TestUUID_Docker(t *testing.T) {
	q.Q("--> starting docker test")
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"uuid": Factory,
		},
	}

	cluster, err := acctest.NewDockerCluster(t.Name(), coreConfig, nil)
	if err != nil {
		t.Fatal(err)
	}

	cores := cluster.ClusterNodes
	client := cores[0].Client

	err = client.Sys().Mount("uuid", &api.MountInput{
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
	q.Q("<-- end docker test")
}
