package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/Nerzal/gocloak/v3"
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

func getClientByClientID(client gocloak.GoCloak, token string, realm string, clientID string) (*gocloak.Client, error) {
	clients, err := client.GetClients(token,
		realm,
		gocloak.GetClientsParams{
			ClientID: clientID,
		})
	if err != nil {
		return nil, err
	}

	for _, fetchedClient := range clients {
		if fetchedClient.ClientID == clientID {
			return fetchedClient, nil
		}
	}
	return nil, nil
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

	token := client.Token
	kcClient := client.Client
	realm := client.Realm

	kcFetchedClient, err := getClientByClientID(kcClient, token.AccessToken, realm, clientId)
	if err != nil {
		return nil, err
	}

	if kcFetchedClient == nil {
		return logical.ErrorResponse(fmt.Sprintf("client %q does not exist on keycloak", clientId)), nil
	}

	creds, err := kcClient.GetClientSecret(token.AccessToken, realm, kcFetchedClient.ID)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": creds.Value,
		},
	}, nil
}
