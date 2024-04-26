# App Service Plan (Consumption Plan for Azure Functions)
resource "azurerm_service_plan" "sp" {
  name                = "serviceplan-gewoscout"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location
  os_type             = "Linux"
  sku_name            = "Y1"
}