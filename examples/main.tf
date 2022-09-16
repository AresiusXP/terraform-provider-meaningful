resource "meaningful_resource_name" "server_name" {
    tenant_id = "766adece-df5e-4735-a192-d80fc644fa8a"
    client_id = "92d60cfa-83c7-4c85-8c42-1cb94720970b"
    client_secret = "clientPassword"
    meaningful_env = "QA"
    resource_type = "Web App"
    deployment_id = "SMXOPX"
    location = "westeurope"
    environment = "Development"
}

output "server_name" {
    value = "${meaningful_resource_name.server_name.name}"
}