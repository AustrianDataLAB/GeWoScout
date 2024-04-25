# Add Application Insights for logs & monitoring
resource "azurerm_application_insights" "ai" {
  name                = "appinsights-gewoscout"
  location            = data.azurerm_resource_group.rg.location
  resource_group_name = data.azurerm_resource_group.rg.name
  application_type    = "other"
}