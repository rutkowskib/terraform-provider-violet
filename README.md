# Terraform Violet Provider

This provider allows managing [Violet](https://violet.io/) webhooks.  

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

### Credentails

To use a provider you need to provide Violet *username*, *password*, *app_id* and *app_secret*. If you dont want to define them in tfvars file to
avoid accidental commiting to repository you can expose them through environmental variables.

```shell
export VIOLET_USERNAME=username
export VIOLET_PASSWORD=password
export VIOLET_APP_SECRET=app_secret
export VIOLET_APP_ID=app_id
```

### Minimal example

The example below is a minimal usage of the provider. It defines a provider and creates a webhook.

```shell
terraform {
    required_providers {
        violet = {
            source = "rutkowskib/violet"
        }
    }
}

provider "violet" {
    username   = var.username
    password   = var.password
    app_id     = var.app_id
    app_secret = var.app_secret
}

resource "violet_webhook" "example" {
    event           = "OFFER_UPDATED"
    remote_endpoint = "https://test.com/"
}
```


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
