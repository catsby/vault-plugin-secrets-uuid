// Package uuid is an example Vault secrets plugin
package uuid

import (
	"context"
	"fmt"
	"strings"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Factory configures and returns uuid backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	var b backend
	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(uuidHelp),
		BackendType: logical.TypeLogical,
		Secrets: []*framework.Secret{
			secretUUID(&b),
		},
	}

	b.Backend.Paths = append(b.Backend.Paths, b.paths()...)

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	// Setup always returns  nil
	_ = b.Backend.Setup(ctx, conf)

	return b, nil
}

// backend wraps the backend framework and adds a map for storing key value pairs
type backend struct {
	*framework.Backend

	store map[string][]byte
}

func (b *backend) paths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: framework.MatchAllRegex("new"),

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
				},
			},
		},
	}
}

func (b *backend) handleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	// output := make(map[string]interface{})
	uuidStr, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("error making UUID: %s", err)
	}
	// output["uuid"] = uuidStr
	// return &logical.Response{
	// 	Data: output,
	// }, nil

	resp := b.Secret(SecretUUIDType).Response(map[string]interface{}{
		"uuid": uuidStr,
	}, nil)
	return resp, nil
}

const uuidHelp = `
The uuid backend is a dummy secrets backend that generates a UUID
`
