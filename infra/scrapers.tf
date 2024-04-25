# App Service Plan (Consumption Plan for Azure Functions)
resource "azurerm_service_plan" "sp_scraper" {
  name                = "serviceplan-gewoscout"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location
  os_type             = "Linux"
  sku_name            = "Y1"
}

# Storage Account name must be globally unique
resource "random_string" "sa_scraper_suffix" {
  length  = 6
  special = false
  upper   = false
}

# Azure Storage Account for the Function App
resource "azurerm_storage_account" "sa_scraper" {
  name                     = "storagescrapergewo${random_string.sa_scraper_suffix.result}"
  resource_group_name      = data.azurerm_resource_group.rg.name
  location                 = data.azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Package the Azure Function's code to zip
data "archive_file" "scraper_zip" {
  type        = "zip"
  source_dir  = "${path.module}/../scrapers"
  output_path = "${path.module}/scrapers-${sha1(join("", [for f in fileset("../scrapers", "**") : filesha1("../scrapers/${f}")]))}.zip"
}

# Ensure Function App name is globally unique
resource "random_string" "fa_scraper_suffix" {
  length  = 6
  special = false
  upper   = false
}

# Azure Function App
resource "azurerm_linux_function_app" "fa_scraper" {
  name                = "funcapp-scraper-gewoscout-${random_string.fa_scraper_suffix.result}"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location

  storage_account_name       = azurerm_storage_account.sa_scraper.name
  storage_account_access_key = azurerm_storage_account.sa_scraper.primary_access_key
  service_plan_id            = azurerm_service_plan.sp_scraper.id

  site_config {
    application_stack {
      python_version = 3.11
    }
    application_insights_connection_string = azurerm_application_insights.ai.connection_string
    application_insights_key               = azurerm_application_insights.ai.instrumentation_key
  }

  app_settings = {
    SCM_DO_BUILD_DURING_DEPLOYMENT = true
    QUEUE_NAME                     = azurerm_storage_queue.queue_scraper_backend.name
  }

  zip_deploy_file = data.archive_file.scraper_zip.output_path
}
