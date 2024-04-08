provider "azurerm" {
  subscription_id = var.infra_subscription_id
  tenant_id       = var.infra_tenant_id
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

terraform {
  backend "azurerm" {
    resource_group_name  = var.state_resource_group_name
    storage_account_name = var.state_storage_account_name
    container_name       = "tfgewoscout"
    key                  = "terraform-base.tfstate"
  }
}

data "azurerm_client_config" "current" {}