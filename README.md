> NOTE: This provider was developed for a specific internal API in 2017. This is deprecated because of changes in the API endpoints, and Terraform version being too old. However, it will provide examples on how to implement API calls with a custom Terraform Provider.

# Meaningful provider for Terraform
Terraform provider to get a name for Azure components following the CT standard. Provider queries Meaningful API and outputs result as string to be used.

## Requirements
* [Terraform](https://www.terraform.io/downloads.html) 0.11.x
* [Go](https://golang.org/doc/install) 1.10.x

##Installing the provider
Follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin). After placing it into your plugins directory, run `terraform init` to initialize it.

##Usage
Provider takes the following parameters, and they're all **required**.

* `tenant_id`: Tenant ID to generate URL in order to request authorization token.
* `client_id`: SPN client ID to request token.
* `client_secret`: SPN secret.
* `meaningful_env`: Meaningful API environment. Valid options are `QA` or `Prod`.
* `resource_type`: Azure component type to be named. Please refer to Meaningful documentation for valid options.
* `deployment_id`: String with 6 characters. First 3 characters are Product, last 3 characters are Application.
* `location`: Azure location where deployment will happen.
* `environment`: Environment type for deployment. Please refer to Meaningful documentation for valid options.


## Example

```ruby
resource "meaningful_resource_name" "server_name" {
    tenant_id = "766adece-df5e-4735-a192-d80fc644fa8a"
    client_id = "92d60cfa-83c7-4c85-8c42-1cb94720970b"
    client_secret = "clientPassword"
    meaningful_env = "QA"
    resource_type = "Web App"
    deployment_id = "ABCDEF"
    location = "westeurope"
    environment = "Development"
}

output "server_name" {
    value = "${meaningful_resource_name.server_name.name}"
}
```

Gives the result:
```
meaningful_resource_name.server_name.name: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

server_name = EUWDABCDEFWAP01
```

## Destroying resource
> When requesting `terraform destroy` it will remove from Meaningful every name generated with the same information; i.e. if you requested a name for several WebApps for Deployment ID ABCDEF, when asking to destroy a specific target with a name generated, it will destroy *every* name for WebApps with ID ABCDEF.

> This is a Meaningful API limitation for deleting names.