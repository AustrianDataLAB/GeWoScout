resource "azurerm_cosmosdb_account" "db_acc" {
  name                = "cosmos-account-gewoscout"
  location            = data.azurerm_resource_group.rg.location
  resource_group_name = data.azurerm_resource_group.rg.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  capabilities {
    name = "EnableServerless"
  }

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = data.azurerm_resource_group.rg.location
    failover_priority = 0
  }

  backup {
    type = "Continuous"
    tier = "Continuous7Days"
  }
}

resource "azurerm_cosmosdb_sql_database" "db" {
  name                = "gewoscout-db"
  resource_group_name = azurerm_cosmosdb_account.db_acc.resource_group_name
  account_name        = azurerm_cosmosdb_account.db_acc.name
}

resource "azurerm_cosmosdb_sql_container" "listings_by_city" {
  name                = "ListingsByCity"
  resource_group_name = azurerm_cosmosdb_account.db_acc.resource_group_name
  account_name        = azurerm_cosmosdb_account.db_acc.name
  database_name       = azurerm_cosmosdb_sql_database.db.name

  partition_key_path = "/_partitionKey"
}

resource "azurerm_cosmosdb_sql_container" "notification_settings_by_city" {
  name                = "NotificationSettingsByCity"
  resource_group_name = azurerm_cosmosdb_account.db_acc.resource_group_name
  account_name        = azurerm_cosmosdb_account.db_acc.name
  database_name       = azurerm_cosmosdb_sql_database.db.name

  partition_key_path = "/_partitionKey"
}

resource "azurerm_cosmosdb_sql_container" "user_data_by_user_id" {
  name                = "UserDataByUserId"
  resource_group_name = azurerm_cosmosdb_account.db_acc.resource_group_name
  account_name        = azurerm_cosmosdb_account.db_acc.name
  database_name       = azurerm_cosmosdb_sql_database.db.name

  partition_key_path = "/_partitionKey"
}