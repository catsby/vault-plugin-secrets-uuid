module github.com/catsby/vault-plugin-secrets-uuid

go 1.14

replace github.com/hashicorp/vault/sdk => /Users/clint/go-src/github.com/hashicorp/vault/sdk

require (
	github.com/hashicorp/go-hclog v0.12.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/terraform-plugin-sdk v1.9.1
	github.com/hashicorp/vault v1.4.0
	github.com/hashicorp/vault/api v1.0.5-0.20200317185738-82f498082f02
	github.com/hashicorp/vault/sdk v0.1.14-0.20200420122737-740110c49f9c
	github.com/y0ssar1an/q v1.0.10
)
