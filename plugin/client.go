package plugin

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v7"
	"github.com/hashicorp/vault/sdk/logical"
)

type Client struct {
	Client gocloak.GoCloak
	Token  *gocloak.JWT
	Realm  string
}

func (b *backend) client(ctx context.Context, s logical.Storage) (*Client, error, error) {
	conf, userErr, intErr := b.readConfig(ctx, s)
	if intErr != nil {
		return nil, nil, intErr
	}
	if userErr != nil {
		return nil, userErr, nil
	}
	if conf == nil {
		return nil, nil, fmt.Errorf("no error received but no config found")
	}

	kcClient := gocloak.NewClient(conf.Url)

	var token *gocloak.JWT
	var err error
	if conf.GrantType == "password" {
		token, err = kcClient.LoginAdmin(ctx, conf.Username, conf.Password, conf.Realm)
	} else if conf.GrantType == "client_credentials" {
		token, err = kcClient.LoginClient(ctx, conf.ClientId, conf.ClientSecret, conf.Realm)
	} else {
		return nil, nil, fmt.Errorf("grant_type %s not implemented", conf.GrantType)
	}

	if err != nil {
		return nil, nil, err
	}
	client := &Client{
		Client: kcClient,
		Token:  token,
		Realm:  conf.Realm,
	}
	return client, nil, nil
}

func (c *Client) getClientSecret(ctx context.Context, clientID string) (string, error) {
	clients, err := c.Client.GetClients(ctx, c.Token.AccessToken,
		c.Realm,
		gocloak.GetClientsParams{
			ClientID: &clientID,
		})
	if err != nil {
		return "", err
	}

	for _, fetchedClient := range clients {
		if *fetchedClient.ClientID == clientID {
			creds, err := c.Client.GetClientSecret(ctx, c.Token.AccessToken, c.Realm, *fetchedClient.ID)
			if err != nil {
				return "", err
			}
			return *creds.Value, nil
		}
	}
	return "", nil
}
