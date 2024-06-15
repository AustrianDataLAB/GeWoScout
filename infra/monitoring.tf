# Add Application Insights for logs & monitoring
resource "azurerm_application_insights" "ai" {
  name                = "appinsights-gewoscout"
  location            = data.azurerm_resource_group.rg.location
  resource_group_name = data.azurerm_resource_group.rg.name
  application_type    = "other"
}

resource "azurerm_portal_dashboard" "main-dashboard" {
  name                = "GewoScout Main Dashboard"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location

  dashboard_properties = templatefile("../dashboards/main.tpl", {
    dashboard_name  = "GewoScout Main Dashboard"
    subscription_id = "e31c37ff-9b82-4f6b-8337-51314cc300ff"
    resource_group  = data.azurerm_resource_group.rg.name
  })
}
