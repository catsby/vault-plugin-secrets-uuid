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
	b := &backend{
		store: make(map[string][]byte),
	}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(uuidHelp),
		BackendType: logical.TypeLogical,
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

			// Fields: map[string]*framework.FieldSchema{
			// 	"new": {
			// 		Type:        framework.TypeString,
			// 		Description: "Generates a new UUID.",
			// 	},
			// },

			Operations: map[logical.Operation]framework.OperationHandler{
				// logical.ReadOperation: &framework.PathOperation{
				// 	Callback: b.handleRead,
				// 	Summary:  "Retrieve the secret from the map.",
				// },
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleWrite,
				},
				// logical.DeleteOperation: &framework.PathOperation{
				// 	Callback: b.handleDelete,
				// 	Summary:  "Deletes the secret at the specified location.",
				// },
			},

			// ExistenceCheck: b.handleExistenceCheck,
		},
	}
}

// TODO return nil/nil here if this is needed
// func (b *backend) handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
// 	out, err := req.Storage.Get(ctx, req.Path)
// 	if err != nil {
// 		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
// 	}

// 	return out != nil, nil
// }

func (b *backend) handleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.ClientToken == "" {
		return nil, fmt.Errorf("client token empty")
	}

	output := make(map[string]interface{})
	uuidStr, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("error making UUID: %s", err)
	}
	output["uuid"] = uuidStr
	return &logical.Response{
		Data: output,
	}, nil
}

const uuidHelp = `
The uuid backend is a dummy secrets backend that generates a UUID
`
