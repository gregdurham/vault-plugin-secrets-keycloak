package plugin

import (
	"context"
	//"strings"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

var b *backend
var storage logical.Storage

func testBackend(tb testing.TB) (*backend, logical.Storage) {
	tb.Helper()

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Factory(context.Background(), config)
	if err != nil {
		tb.Fatal(err)
	}
	return b.(*backend), config.StorageView
}

func TestBackend(t *testing.T) {
	b, storage = testBackend(t)

	t.Run("write config", WriteConfig)
	t.Run("read config", ReadConfig)
}

func WriteConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   storage,
		Data: map[string]interface{}{
			"url":      "http://localhost:8080",
			"username": "admin",
			"password": "password",
			"realm":    "master",
		},
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatal("expected no response because Vault generally doesn't return it for posts")
	}
}

func ReadConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   storage,
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(err)
	}

	expectedKeys := map[string]bool{
		"url":   true,
		"realm": true,
	}

	for k := range resp.Data {
		if !expectedKeys[k] {
			t.Fatalf("%q exists in response, but should not be", fmt.Sprintf("%s", k))
		}
	}

	if resp.Data["url"] != `http://localhost:8080` {
		t.Fatalf("url to be \"+http://localhost:8080+\" but received %q", fmt.Sprintf("%s", resp.Data["url"]))
	}

	if resp.Data["realm"] != `master` {
		t.Fatalf("realm to be \"+master+\" but received %q", fmt.Sprintf("%s", resp.Data["realm"]))
	}
}
