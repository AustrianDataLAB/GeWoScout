{
  "lenses": {
    "0": {
      "order": 0,
      "parts": {
        "0": {
          "position": {
            "x": 0,
            "y": 0,
            "colSpan": 3,
            "rowSpan": 5
          },
          "metadata": {
            "inputs": [],
            "type": "Extension/HubsExtension/PartType/MarkdownPart",
            "settings": {
              "content": {
                "content": "# Scrapers\nThe scrapers are implemented as an Azure Function App, and scrape Genossenschafts websites according to a cron-schedule. After successful data-extraction, information about appartment listings gets written to the `scraper-result-queue` queue.",
                "title": "",
                "subtitle": "",
                "markdownSource": 1,
                "markdownUri": ""
              }
            }
          }
        },
        "1": {
          "position": {
            "x": 3,
            "y": 0,
            "colSpan": 10,
            "rowSpan": 1
          },
          "metadata": {
            "inputs": [],
            "type": "Extension/HubsExtension/PartType/MarkdownPart",
            "settings": {
              "content": {
                "content": "## BWSG Scraper",
                "title": "",
                "subtitle": "",
                "markdownSource": 1,
                "markdownUri": ""
              }
            }
          }
        },
        "2": {
          "position": {
            "x": 3,
            "y": 1,
            "colSpan": 4,
            "rowSpan": 4
          },
          "metadata": {
            "inputs": [
              {
                "name": "sharedTimeRange",
                "isOptional": true
              },
              {
                "name": "options",
                "value": {
                  "chart": {
                    "title": "Successful Execution Count - BWSG",
                    "metrics": [
                      {
                        "name": "bwsg_scraper Successes",
                        "resourceMetadata": {
                          "id": "/subscriptions/${subscription_id}/resourceGroups/${resource_group}/providers/microsoft.insights/components/appinsights-gewoscout"
                        },
                        "aggregationType": 1
                      }
                    ],
                    "visualization": {
                      "disablePinning": true
                    }
                  }
                },
                "isOptional": true
              }
            ],
            "type": "Extension/HubsExtension/PartType/MonitorChartPart",
            "settings": {
              "content": {
                "options": {
                  "chart": {
                    "metrics": [
                      {
                        "resourceMetadata": {
                          "id": "/subscriptions/${subscription_id}/resourceGroups/${resource_group}/providers/microsoft.insights/components/appinsights-gewoscout"
                        },
                        "name": "customMetrics/bwsg_scraper Successes",
                        "aggregationType": 1,
                        "namespace": "microsoft.insights/components/kusto",
                        "metricVisualization": {
                          "displayName": "bwsg_scraper Successes"
                        }
                      }
                    ],
                    "title": "Successful Execution Count",
                    "titleKind": 2,
                    "visualization": {
                      "chartType": 2,
                      "legendVisualization": {
                        "isVisible": true,
                        "position": 2,
                        "hideHoverCard": false,
                        "hideLabelNames": true
                      },
                      "axisVisualization": {
                        "x": {
                          "isVisible": true,
                          "axisType": 2
                        },
                        "y": {
                          "isVisible": true,
                          "axisType": 1
                        }
                      },
                      "disablePinning": true
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "metadata": {
    "model": {
      "timeRange": {
        "value": {
          "relative": {
            "duration": 24,
            "timeUnit": 1
          }
        },
        "type": "MsPortalFx.Composition.Configuration.ValueTypes.TimeRange"
      },
      "filterLocale": {
        "value": "en-us"
      },
      "filters": {
        "value": {
          "MsPortalFx_TimeRange": {
            "model": {
              "format": "local",
              "granularity": "auto",
              "relative": "1h"
            },
            "displayCache": {
              "name": "Local Time",
              "value": "Past hour"
            },
            "filteredPartIds": [
              "StartboardPart-MonitorChartPart-86a97885-31ad-4d31-90ba-35d0b493ff34"
            ]
          }
        }
      }
    }
  }
}