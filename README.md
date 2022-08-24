# Commercelayer Terraform Provider

The Commercelayer terraform provider allows you to configure your [Commercelayer shops](https://commercelayer.io/) with
infrastructure-as-code principles.

## Usage

Add the provider to your terraform project

```hcl
terraform {
  required_providers {
    commercelayer = {
      version = ">= 0.0.1"
      source  = "incentro-dc/commercelayer"
    }
  }
}

provider "commercelayer" {
  client_id     = "<client_id>"
  client_secret = "<client_secret>"
  api_endpoint  = "<api_endpoint>"
  auth_endpoint = "<auth_endpoint"
}
```

## Development

### Requirements

In order to build from the source code, you must have the following set up in your development environment.

- [Go >= 1.17](https://golang.org/doc/install)
- [Make](https://www.gnu.org/software/make/)
- [Terraform >= 1.0.0](https://www.terraform.io/downloads.html)

There is also a dependency on another internal
project, [which provides the SDK used](https://github.com/incentro-dc/go-commercelayer-sdk).

### Running

Build the binary with `make`. Note that this will also import any required dependencies and generate any code or
documentation necessary.

    make build

This will produce the project binary. Note that by default the `go build` process will check your environment and build
the binary (using the project name) accordingly. If you want to change this check out the build options `go help build`.

Now you can run the binary

    ./terraform-provider-commercelayer

This will however only tell you that the project needs to run as a plugin to terraform. To this end we can also provide
a parameter to the binary to tell it to run in development mode

    ./terraform-provider-commercelayer -debug

This will provide an environment variable that can be loaded when initializing and running terraform

    export TF_REATTACH_PROVIDERS='<provider data>'
    terraform init
    terraform apply

### Testing and cleaning

Run the tests to check for any issues

```
go test ./...
```

Run formatting to clean up the code (you might need to run this several times to make sure all issues have been handled)

```
go fmt ./...
```

Make sure to also clean up the mod file

```
go mod tidy
```

## Examples

See the [examples folder](./examples) for some examples of terraform code