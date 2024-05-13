locals {
  backend_path = "${path.module}/../backend"
}

# Random suffix for the storage account
resource "random_string" "sa_backend_suffix" {
  length  = 6
  special = false
  upper   = false
}

# Azure Storage Account for the Function App
resource "azurerm_storage_account" "sa_backend" {
  name                     = "sabackend${random_string.sa_backend_suffix.result}"
  resource_group_name      = data.azurerm_resource_group.rg.name
  location                 = data.azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Random suffix for the function
resource "random_string" "fa_backend_suffix" {
  length  = 6
  special = false
  upper   = false
}

# Azure Function App
resource "azurerm_linux_function_app" "fa_backend" {
  name                = "funcapp-backend-gewoscout-${random_string.fa_backend_suffix.result}"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location

  storage_account_name       = azurerm_storage_account.sa_backend.name
  storage_account_access_key = azurerm_storage_account.sa_backend.primary_access_key
  service_plan_id            = azurerm_service_plan.sp.id

  site_config {
    application_stack {
      use_custom_runtime = true
    }
    application_insights_connection_string = azurerm_application_insights.ai.connection_string
    application_insights_key               = azurerm_application_insights.ai.instrumentation_key
  }

  app_settings = {
    FUNCTIONS_WORKER_RUNTIME = "custom"
    QUEUE_NAME               = azurerm_storage_queue.queue_scraper_backend.name
    COSMOS_DB_CONNECTION     = azurerm_cosmosdb_account.db_acc.primary_sql_connection_string
  }
  
  zip_deploy_file = data.archive_file.backend_zip.output_path
}


# Build the backend
resource "null_resource" "backend_build" {
  # Using triggers to force execution on every apply
  triggers = {
    always_run = timestamp()
  }

  depends_on = [azurerm_linux_function_app.fa_backend]

  provisioner "local-exec" {
    working_dir = local.backend_path
    command     = "make build-linux-minimal"
  }
}

# Package the Azure Function's code to zip
data "archive_file" "backend_zip" {
  type        = "zip"
  source_dir  = "${path.module}/../backend"
  output_path = "${path.module}/be-${sha1(join("", [for f in fileset("../backend", "**") : filesha1("../backend/${f}")]))}.zip"
}
