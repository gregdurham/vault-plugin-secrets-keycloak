package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCredential(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("clientId"),

		Fields: map[string]*framework.FieldSchema{
			"clientId": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The clientId for which to retrieve the client credentials",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCredentialRead,
		},
	}
}

func (b *backend) pathCredentialRead(ctx context.Context, req *logical.Request,
	d *framework.FieldData) (*logical.Response, error) {

	clientId := d.Get("clientId").(string)

	client, userErr, intErr := b.client(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}

	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}

	secret, err := client.getClientSecret(ctx, clientId)

	if err != nil {
		return nil, err
	}

	if secret == "" {
		return logical.ErrorResponse(fmt.Sprintf("client %q does not exist on keycloak", clientId)), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": secret,
		},
	}, nil
}
