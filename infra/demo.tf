terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.0"
    }
  }
  backend "azurerm" {
    use_azuread_auth = true
  }
  required_version = ">= 0.13"
}

variable "subscription_id" {
  type = string
}

variable "tenant_id" {
  type = string
}

provider "azurerm" {
  subscription_id   = var.subscription_id
  tenant_id         = var.tenant_id
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

# TODO read the rg created by C for GeWoScout
data "azurerm_resource_group" "rgruntime" {
  name = "rg-service-gepip35"
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
