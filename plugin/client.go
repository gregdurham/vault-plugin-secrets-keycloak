package plugin

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v3"
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
	token, err := kcClient.LoginAdmin(conf.Username, conf.Password, conf.Realm)

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
