// sources: https://learn.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-static-website-terraform?tabs=azure-cli

variable "frontend_path" {
  type    = string
  default = "../frontend/vue-gewoscout"
}

# Generate a random value for the storage account name
resource "random_string" "sa_frontend_suffix" {
  length  = 6
  upper   = false
  numeric = false
  special = false
}

# Prepare a storage account for the frontend application
resource "azurerm_storage_account" "frontend_storage_account" {
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location

  name = "safrontend${random_string.sa_frontend_suffix.result}"

  account_tier             = "Standard"
  account_replication_type = "LRS"

  static_website {
    index_document = "index.html"
  }
}

# Build the frontend
resource "null_resource" "frontend_build" {
  # Using triggers to force execution on every apply
  triggers = {
    always_run = timestamp()
  }

  depends_on = [azurerm_storage_account.frontend_storage_account]

  provisioner "local-exec" {
    working_dir = var.frontend_path
    command     = <<-EOT
      npm install && npm run build
    EOT
  }
}

# Upload frontend files to the newly created storage account
resource "null_resource" "frontend_upload" {
  # Using triggers to force execution on every apply
  triggers = {
    always_run = timestamp()
  }

  depends_on = [null_resource.frontend_build]

  provisioner "local-exec" {
    working_dir = var.frontend_path
    command     = <<EOT
      az login --service-principal -u ${var.arm_client_id} -p ${var.arm_client_secret} --tenant ${var.arm_tenant_id} && az storage blob upload-batch --overwrite -s ./dist -d $web --account-name ${azurerm_storage_account.frontend_storage_account.name}
    EOT
  }
}

output "frontend_hostname" {
  description = "Frontend hostname"
  value       = azurerm_storage_account.frontend_storage_account.primary_web_host
}