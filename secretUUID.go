package uuid

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// SecretUUIDType is the name used to identify this type
const SecretUUIDType = "uuid"

func secretUUID(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretUUIDType,
		Fields: map[string]*framework.FieldSchema{
			"uuid": {
				Type: framework.TypeString,
				Description: `The PEM-encoded concatenated certificate and
issuing certificate authority`,
			},
		},

		Revoke: b.secretUUIDRevoke,
	}
}

func (b *backend) secretUUIDRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}
