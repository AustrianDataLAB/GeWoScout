# Storage Queue - connects scrapers with backend
resource "azurerm_storage_queue" "queue_scraper_backend" {
  name                 = "queue-scraper-backend-gewoscout"
  storage_account_name = azurerm_storage_account.sa_scraper.name
}
