resource "random_string" "sa_queue_suffix" {
  length  = 6
  special = false
  upper   = false
}

resource "azurerm_storage_account" "sa_queue" {
  name                     = "saqueue${random_string.sa_queue_suffix.result}"
  resource_group_name      = data.azurerm_resource_group.rg.name
  location                 = data.azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# Storage Queue - connects scrapers with backend
resource "azurerm_storage_queue" "queue_scraper_backend" {
  name                 = "scraping-input-node"
  storage_account_name = azurerm_storage_account.sa_queue.name
}