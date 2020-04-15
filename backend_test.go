package uuid

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/acctest"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/y0ssar1an/q"
)

var helper *acctest.Helper

// RunningAsPlugin returns true if it detects the usual Terraform plugin
// detection environment variables, suggesting that the current process is
// being launched as a plugin server.
// TODO: move to acctest package
func RunningAsPlugin() bool {
	magicCookieKey :=
		"VAULT_BACKEND_PLUGIN"
	magicCookieValue :=
		"6669da05-b1c8-4f49-97d9-c8e5bed98e20"

	return os.Getenv(magicCookieKey) == magicCookieValue
}

func TestMain(m *testing.M) {
	// run acc tests
	if os.Getenv("VAULT_ACC") == "1" {
		absPluginExecPath, _ := filepath.Abs(os.Args[0])
		pluginName := path.Base(absPluginExecPath)
		os.Link(absPluginExecPath, path.Join("/Users/clint/Desktop/plugins", pluginName))
		// setup docker, send src and name
		// run tests
		coreConfig := &vault.CoreConfig{
			LogicalBackends: map[string]logical.Factory{
				"uuid": Factory,
			},
		}
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		wd = path.Join(wd, "vault/plugins/uuid")
		cmd := exec.Command("go", "build", "-o", "./vault/plugins/uuid", "/Users/clint/go-src/github.com/catsby/vault-plugin-secrets-uuid/cmd/uuid/main.go")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		// cluster, err := acctest.NewDockerCluster(t.Name(), coreConfig, nil)
		// dOpts := &acctest.DockerClusterOptions{PluginTestBin: absPluginExecPath}
		//TODO: cleanup working dir
		dOpts := &acctest.DockerClusterOptions{PluginTestBin: wd}
		cluster, err := acctest.NewDockerCluster("test-uuid", coreConfig, dOpts)
		if err != nil {
			panic(err)
		}

		cores := cluster.ClusterNodes
		client := cores[0].Client
		// calculate sha256 of binary/vault/plugins/uuid
		pPath := "/Users/clint/go-src/github.com/catsby/vault-plugin-secrets-uuid/vault/plugins/uuid"

		f, err := os.Open(pPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			panic(err)
		}
		sha256value := fmt.Sprintf("%x", h.Sum(nil))
		q.Q("sha256:", sha256value)

		helper = &acctest.Helper{
			Client: client,
		}
		// use client to mount plugin

		err = client.Sys().RegisterPlugin(&api.RegisterPluginInput{
			Name:    "uuid",
			Type:    consts.PluginTypeSecrets,
			Command: "uuid",
			SHA256:  sha256value,
		})
		if err != nil {
			panic(err)
		}

		// run tests
		os.Exit(m.Run())
	} else {
		// run normal test
		os.Exit(m.Run())
	}
}

func TestAccUUID_Docker(t *testing.T) {
	if os.Getenv("VAULT_ACC") == "" {
		t.Log("VAULT_ACC is not set")
		t.SkipNow()
	}
	q.Q("--> starting docker test")
	if helper == nil {
		t.Fatal("nil helper")
	}
	client := helper.Client
	// s, lErr := client.Logical().Read("/sys/mounts")
	// if lErr != nil {
	// 	q.Q("list err:", lErr)
	// } else {
	// 	q.Q("secrets list not error")
	// }
	// q.Q("secret list s:", s)

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
