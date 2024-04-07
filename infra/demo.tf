terraform {
  backend "azurerm" {
    resource_group_name  = "<resourcegroupname>"
    storage_account_name = "<storageaccount-name>"
    container_name       = "<container-name"
    key                  = "<Key Value>"
  }
}
 
provider "azurerm" {
  # The "feature" block is required for AzureRM provider > 2.x.
  version = ">= 3.50.0"
  features {}
}
 
data "azurerm_client_config" "current" {}
 
#Create Resource Group
resource "azurerm_resource_group" "rgsa" {
  name     = "<resourcegroup1-name>"
  location = "eastus2"
}
 
#Create Virtual Network
resource "azurerm_virtual_network" "vnet" {
  name                = "rgsa0102-vnet"
  address_space       = ["<CIDR Range>"]
  location            = "eastus2"
  resource_group_name = azurerm_resource_group.rgsa.name
}
 
# Create Subnet
resource "azurerm_subnet" "subnet" {
  name                 = "subnet"
  resource_group_name  = azurerm_resource_group.rgsa.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefix       = "<CIDR Range>"
}