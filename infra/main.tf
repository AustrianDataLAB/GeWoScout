# TODO read the rg created by C for GeWoScout
data "azurerm_resource_group" "rgruntime" {
  name = var.infra_resource_group_name
}

// storage account contains the filesystem of our function
resource "azurerm_storage_account" "storage" {
  name                     = "gewoscoutdevstorage"
  resource_group_name      = data.azurerm_resource_group.rgruntime.name
  location                 = data.azurerm_resource_group.rgruntime.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

// defines resources available to function: Y1 has free tier
resource "azurerm_service_plan" "serviceplan" {
  name                = "gewoscout-dev-serviceplan"
  resource_group_name = data.azurerm_resource_group.rgruntime.name
  location            = data.azurerm_resource_group.rgruntime.location
  os_type             = "Linux"
  sku_name            = "Y1"
}

// function
resource "azurerm_linux_function_app" "function" {
  name                = "gewoscout-dev-function-app"
  resource_group_name = data.azurerm_resource_group.rgruntime.name
  location            = data.azurerm_resource_group.rgruntime.location

  storage_account_name       = azurerm_storage_account.storage.name
  storage_account_access_key = azurerm_storage_account.storage.primary_access_key
  service_plan_id            = azurerm_service_plan.serviceplan.id

  site_config {
    application_stack {
      python_version = "3.9"
    }
  }
}
