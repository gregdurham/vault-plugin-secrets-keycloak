package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type keycloakConfig struct {
	Url          string `json:"url"`
	GrantType    string `json:"grant_type"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Realm        string `json:"realm"`
}

func pathConfigConnection(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/connection",

		Fields: map[string]*framework.FieldSchema{
			"url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Keycloak Server URL",
			},
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Keycloak Admin username",
			},
			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Keycloak Admin password",
			},
			"client_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Keycloak client_id",
			},
			"client_secret": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Keycloak client_secret",
			},
			"realm": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Keycloak realm",
			},
			"grant_type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "grant_type used for authenticating to keycloak",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigConnectionRead,
			logical.UpdateOperation: b.pathConfigConnectionWrite,
		},
	}
}

func (b *backend) readConfig(ctx context.Context, storage logical.Storage) (*keycloakConfig, error, error) {
	entry, err := storage.Get(ctx, "config/connection")
	if err != nil {
		return nil, nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("Access credentials for the backend have not been configured; please configure them at the '/config' endpoint"), nil
	}

	conf := &keycloakConfig{}
	if err := entry.DecodeJSON(conf); err != nil {
		return nil, nil, fmt.Errorf("Error reading vault config")
	}

	return conf, nil, nil
}

func (b *backend) pathConfigConnectionRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, userErr, intErr := b.readConfig(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}
	if conf == nil {
		return nil, fmt.Errorf("no error reported but vault access configuration not found")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"url":   conf.Url,
			"realm": conf.Realm,
		},
	}, nil

}

func (b *backend) pathConfigConnectionWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	url := data.Get("url").(string)
	username := data.Get("username").(string)
	password := data.Get("password").(string)
	clientId := data.Get("client_id").(string)
	clientSecret := data.Get("client_secret").(string)
	realm := data.Get("realm").(string)
	grantType := data.Get("grant_type").(string)

	if url == "" {
		return logical.ErrorResponse("url parameter must be supplied"), nil
	}

	if realm == "" {
		return logical.ErrorResponse("realm parameter must be supplied"), nil
	}

	if grantType == "" {
		grantType = "password"
	}

	switch grantType {
	case "password":
		if username == "" {
			return logical.ErrorResponse("Username parameter must be supplied when using password grant_type"), nil
		}

		if password == "" {
			return logical.ErrorResponse("Password parameter must be supplied when using password grant_type"), nil
		}
	case "client_credentials":
		if clientId == "" {
			return logical.ErrorResponse("clientId parameter must be supplied when using client_credentials grant_type"), nil
		}

		if clientSecret == "" {
			return logical.ErrorResponse("clientSecret parameter must be supplied when using client_credentials grant_type"), nil
		}
	default:
		return logical.ErrorResponse("Invalid grant_type supplied, available options are password and client_credentials"), nil
	}

	entry, err := logical.StorageEntryJSON("config/connection", keycloakConfig{
		Url:          url,
		GrantType:    grantType,
		Username:     username,
		Password:     password,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Realm:        realm,
	})

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}
