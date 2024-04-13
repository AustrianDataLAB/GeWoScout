# App Service Plan (Consumption Plan for Azure Functions)
resource "azurerm_service_plan" "service_plan" {
  name                = "dev-serviceplan-gewoscout"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location
  os_type             = "Linux"
  sku_name            = "Y1"
}

# Azure Storage Account for the Function App
resource "azurerm_storage_account" "storage_func" {
  name                     = "devfuncstoragegewoscout"
  resource_group_name      = data.azurerm_resource_group.rg.name
  location                 = data.azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Package the Azure Function Python code into a zip file
data "archive_file" "scraper_zip" {
  type        = "zip"
  source_dir  = "${path.module}/../test_http_function"
  output_path = "${path.module}/function.zip"
}

# Azure Function App
resource "azurerm_linux_function_app" "func_scraper" {
  name                = "dev-funcapp-gewoscout"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location

  storage_account_name       = azurerm_storage_account.storage_func.name
  storage_account_access_key = azurerm_storage_account.storage_func.primary_access_key
  service_plan_id            = azurerm_service_plan.service_plan.id

  // TODO remove for scraping function
  https_only = true

  site_config {
    application_stack {
      python_version = "3.11"
    }
  }

  /*
  app_settings = {
    # Dynamically set the Cosmos DB connection string
    "CosmosDBConnectionString" = azurerm_cosmosdb_account.gewoscout_cosmosdb.primary_connection_string
  }
  */

  zip_deploy_file = data.archive_file.scraper_zip.output_path
}

// todo add application insights