# Vault Plugin: Keycloak Secrets Backend

Vault secrets plugin for Keycloak. Manage keycloak client secret credentials via Hashicorp Vault

## Quick Links
    - Vault Website: https://www.vaultproject.io
    - Keycloak Website: https://www.keycloak.org

## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

## Usage & API

#### Configuration

| Method   | Path                                       | Produces                 |
| :------- | :----------------------------------------- | :----------------------- |
| `POST`   | `/keycloak/config/connection`              | `200 (application/json)` |

#### Parameters
- `url` `(string)` - Specifies the url to the keycloak server, format should be: 
  http://127.0.0.1:8080
- `username` `(string)` - Specifies the username to authenticate to keycloak with
- `password` `(string)` - Specifies the password for the user to authenticate to 
  keycloak with
- `realm` `(string)` - Specifies the realm in which the client(s) which credentials 
  will be requested for as well as the user exists 

| Method   | Path                                       | Produces                 |
| :------- | :----------------------------------------- | :----------------------- |
| `GET`   | `/keycloak/creds/:name`                     | `200 (application/json)` |

#### Parameters
- `name` `(string)` - Specifies the name of the keycloak client to request the credentials for

#### CLI

```sh
$ vault secrets enable keycloak
Success! Enabled the keycloak secrets engine at: keycloak/
```

```sh
$ vault write keycloak/config/connection url=http://localhost:8080 username=admin password=password realm=master
```

```sh
$ vault read keycloak/creds/clientName
```

## Developing

If you wish to work on this plugin, you'll first need
[Go](https://www.golang.org) installed on your machine
(version 1.10+ is *required*).

For local dev first make sure Go is properly installed, including
setting up a [GOPATH](https://golang.org/doc/code.html#GOPATH).
Next, clone this repository into
`$GOPATH/src/github.com/gregdurham/vault-plugin-secrets-keycloak`.
You can then download any required build tools by bootstrapping your
environment:

```sh
$ make bootstrap
```

To compile a development version of this plugin, run `make` or `make dev`.
This will put the plugin binary in the `bin` and `$GOPATH/bin` folders. `dev`
mode will only generate the binary for your platform and is faster:

```sh
$ make
$ make dev
```

Put the plugin binary into a location of your choice. This directory
will be specified as the [`plugin_directory`](https://www.vaultproject.io/docs/configuration/index.html#plugin_directory)
in the Vault config used to start the server.

```json
...
plugin_directory = "path/to/plugin/directory"
...
```

Start a Vault server with this config file:
```sh
$ vault server -config=path/to/config.json ...
...
```

Once the server is started, register the plugin in the Vault server's [plugin catalog](https://www.vaultproject.io/docs/internals/plugins.html#plugin-catalog):

```sh
$ vault write sys/plugins/catalog/secret/keycloak \
        sha_256=<expected SHA256 Hex value of the plugin binary> \
        command="vault-plugin-secrets-keycloak"
...
Success! Data written to: sys/plugins/catalog/secret/keycloak
```

Note you should generate a new sha256 checksum if you have made changes
to the plugin. Example using openssl:

```sh
openssl dgst -sha256 $GOPATH/vault-plugin-secrets-keycloak
...
SHA256(.../go/bin/vault-plugin-secrets-keycloak)= 896c13c0f5305daed381952a128322e02bc28a57d0c862a78cbc2ea66e8c6fa1
```

Enable the secrets plugin backend using the secrets enable plugin command:

```sh
$ vault secrets enable keycloak
...

Successfully enabled 'plugin' at 'keycloak'!
```

#### Tests

If you are developing this plugin and want to verify it is still
functioning (and you haven't broken anything else), we recommend
running the tests.

To run the tests, invoke `make test`:

```sh
$ make test
```

You can also specify a `TESTARGS` variable to filter tests like so:

```sh
$ make test TESTARGS='--run=TestConfig'
```