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
                "markdownSource": 1,
                "markdownUri": "",
                "subtitle": "",
                "title": ""
              }
            }
          }
        },
        "1": {
          "position": {
            "x": 3,
            "y": 0,
            "colSpan": 12,
            "rowSpan": 1
          },
          "metadata": {
            "inputs": [],
            "type": "Extension/HubsExtension/PartType/MarkdownPart",
            "settings": {
              "content": {
                "content": "## BWSG Scraper",
                "markdownSource": 1,
                "markdownUri": "",
                "subtitle": "",
                "title": ""
              }
            }
          }
        },
        "2": {
          "position": {
            "x": 3,
            "y": 1,
            "colSpan": 6,
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
                    "metrics": [
                      {
                        "aggregationType": 1,
                        "name": "bwsg_scraper Successes",
                        "resourceMetadata": {
                          "id": "/subscriptions/${subscription_id}/resourceGroups/${resource_group}/providers/microsoft.insights/components/${app_insights_name}"
                        }
                      }
                    ],
                    "title": "Successful Execution Count - BWSG",
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
                          "id": "/subscriptions/${subscription_id}/resourceGroups/${resource_group}/providers/microsoft.insights/components/${app_insights_name}"
                        },
                        "name": "customMetrics/bwsg_scraper Successes",
                        "aggregationType": 1,
                        "namespace": "microsoft.insights/components/kusto",
                        "metricVisualization": {
                          "displayName": "bwsg_scraper Successes"
                        }
                      }
                    ],
                    "title": "Successful execution count - bwsg_scraper",
                    "titleKind": 2,
                    "visualization": {
                      "chartType": 1,
                      "legendVisualization": {
                        "hideHoverCard": false,
                        "hideLabelNames": true,
                        "isVisible": true,
                        "position": 2
                      },
                      "axisVisualization": {
                        "x": {
                          "axisType": 2,
                          "isVisible": true
                        },
                        "y": {
                          "axisType": 1,
                          "isVisible": true
                        }
                      },
                      "disablePinning": true
                    }
                  }
                }
              }
            }
          }
        },
        "3": {
          "position": {
            "x": 9,
            "y": 1,
            "colSpan": 6,
            "rowSpan": 4
          },
          "metadata": {
            "inputs": [
              {
                "name": "options",
                "isOptional": true
              },
              {
                "name": "sharedTimeRange",
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
                          "id": "/subscriptions/${subscription_id}/resourceGroups/${resource_group}/providers/Microsoft.Insights/components/${app_insights_name}"
                        },
                        "name": "customMetrics/bwsg_scraper Failures",
                        "aggregationType": 1,
                        "namespace": "microsoft.insights/components/kusto",
                        "metricVisualization": {
                          "displayName": "bwsg_scraper Failures"
                        }
                      }
                    ],
                    "title": "Failed execution count - bwsg_scraper",
                    "titleKind": 2,
                    "visualization": {
                      "chartType": 1,
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
        },
        "4": {
          "position": {
            "x": 3,
            "y": 5,
            "colSpan": 12,
            "rowSpan": 1
          },
          "metadata": {
            "inputs": [],
            "type": "Extension/HubsExtension/PartType/MarkdownPart",
            "settings": {
              "content": {
                "content": "## ÖVW Scraper",
                "title": "",
                "subtitle": "",
                "markdownSource": 1,
                "markdownUri": ""
              }
            }
          }
        },
        "5": {
          "position": {
            "x": 3,
            "y": 6,
            "colSpan": 6,
            "rowSpan": 4
          },
          "metadata": {
            "inputs": [
              {
                "name": "options",
                "isOptional": true
              },
              {
                "name": "sharedTimeRange",
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
                          "id": "/subscriptions/${subscription_id}/resourceGroups/${resource_group}/providers/Microsoft.Insights/components/${app_insights_name}"
                        },
                        "name": "customMetrics/oevw_scraper Successes",
                        "aggregationType": 1,
                        "namespace": "microsoft.insights/components/kusto",
                        "metricVisualization": {
                          "displayName": "oevw_scraper Successes"
                        }
                      }
                    ],
                    "title": "Successful execution count - oevw_scraper",
                    "titleKind": 2,
                    "visualization": {
                      "chartType": 1,
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
        },
        "6": {
          "position": {
            "x": 9,
            "y": 6,
            "colSpan": 6,
            "rowSpan": 4
          },
          "metadata": {
            "inputs": [
              {
                "name": "options",
                "isOptional": true
              },
              {
                "name": "sharedTimeRange",
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
                          "id": "/subscriptions/${subscription_id}/resourceGroups/${resource_group}/providers/Microsoft.Insights/components/${app_insights_name}"
                        },
                        "name": "customMetrics/oevw_scraper Failures",
                        "aggregationType": 1,
                        "namespace": "microsoft.insights/components/kusto",
                        "metricVisualization": {
                          "displayName": "oevw_scraper Failures"
                        }
                      }
                    ],
                    "title": "Failed execution count - oevw_scraper",
                    "titleKind": 2,
                    "visualization": {
                      "chartType": 1,
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
        },
        "7": {
          "position": {
            "x": 3,
            "y": 10,
            "colSpan": 12,
            "rowSpan": 1
          },
          "metadata": {
            "inputs": [],
            "type": "Extension/HubsExtension/PartType/MarkdownPart",
            "settings": {
              "content": {
                "content": "## WBV-GPA Scraper",
                "title": "",
                "subtitle": "",
                "markdownSource": 1,
                "markdownUri": ""
              }
            }
          }
        },
        "8": {
          "position": {
            "x": 3,
            "y": 11,
            "colSpan": 6,
            "rowSpan": 4
          },
          "metadata": {
            "inputs": [
              {
                "name": "options",
                "isOptional": true
              },
              {
                "name": "sharedTimeRange",
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
                          "id": "/subscriptions/${subscription_id}/resourceGroups/${resource_group}/providers/Microsoft.Insights/components/${app_insights_name}"
                        },
                        "name": "customMetrics/wbv_gpa_scraper Successes",
                        "aggregationType": 1,
                        "namespace": "microsoft.insights/components/kusto",
                        "metricVisualization": {
                          "displayName": "wbv_gpa_scraper Successes"
                        }
                      }
                    ],
                    "title": "Successful execution count - wbv_gpa_scraper",
                    "titleKind": 2,
                    "visualization": {
                      "chartType": 1,
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
        },
        "9": {
          "position": {
            "x": 9,
            "y": 11,
            "colSpan": 6,
            "rowSpan": 4
          },
          "metadata": {
            "inputs": [
              {
                "name": "options",
                "isOptional": true
              },
              {
                "name": "sharedTimeRange",
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
                          "id": "/subscriptions/${subscription_id}/resourceGroups/${resource_group}/providers/Microsoft.Insights/components/${app_insights_name}"
                        },
                        "name": "customMetrics/wbv_gpa_scraper Failures",
                        "aggregationType": 1,
                        "namespace": "microsoft.insights/components/kusto",
                        "metricVisualization": {
                          "displayName": "wbv_gpa_scraper Failures"
                        }
                      }
                    ],
                    "title": "Failed execution count - wbv_gpa_scraper",
                    "titleKind": 2,
                    "visualization": {
                      "chartType": 1,
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
        },
        "10": {
          "position": {
            "x": 0,
            "y": 15,
            "colSpan": 3,
            "rowSpan": 4
          },
          "metadata": {
            "inputs": [],
            "type": "Extension/HubsExtension/PartType/MarkdownPart",
            "settings": {
              "content": {
                "content": "# Backend\nThe backend exposes APIs to \n- frontend for retrieving listing data\n- swagger\n- healthchecks\n\nIt also contains endpoints that get triggered by\n- new messages in the scraper queue\n- cosmos events",
                "title": "",
                "subtitle": "",
                "markdownSource": 1,
                "markdownUri": ""
              }
            }
          }
        },
        "11": {
          "position": {
            "x": 3,
            "y": 15,
            "colSpan": 6,
            "rowSpan": 4
          },
          "metadata": {
            "inputs": [
              {
                "name": "options",
                "isOptional": true
              },
              {
                "name": "sharedTimeRange",
                "isOptional": true
              }
            ],
            "type": "Extension/HubsExtension/PartType/MonitorChartPart",
            "settings": {
              "content": {
                "options": {
                  "chart": {
                    "metrics": [],
                    "title": "Metrics chart",
                    "titleKind": 1,
                    "visualization": {
                      "disablePinning": true
                    }
                  }
                }
              }
            }
          }
        },
        "12": {
          "position": {
            "x": 9,
            "y": 15,
            "colSpan": 6,
            "rowSpan": 4
          },
          "metadata": {
            "inputs": [
              {
                "name": "options",
                "isOptional": true
              },
              {
                "name": "sharedTimeRange",
                "isOptional": true
              }
            ],
            "type": "Extension/HubsExtension/PartType/MonitorChartPart",
            "settings": {
              "content": {
                "options": {
                  "chart": {
                    "metrics": [],
                    "title": "Metrics chart",
                    "titleKind": 1,
                    "visualization": {
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
              "StartboardPart-MonitorChartPart-c16f8462-0903-47eb-9700-b8292a7c808a",
              "StartboardPart-MonitorChartPart-c16f8462-0903-47eb-9700-b8292a7c808c",
              "StartboardPart-MonitorChartPart-c16f8462-0903-47eb-9700-b8292a7c8090",
              "StartboardPart-MonitorChartPart-c16f8462-0903-47eb-9700-b8292a7c8092",
              "StartboardPart-MonitorChartPart-c16f8462-0903-47eb-9700-b8292a7c8096",
              "StartboardPart-MonitorChartPart-c16f8462-0903-47eb-9700-b8292a7c8098",
              "StartboardPart-MonitorChartPart-c16f8462-0903-47eb-9700-b8292a7c846c",
              "StartboardPart-MonitorChartPart-c16f8462-0903-47eb-9700-b8292a7c8478"
            ]
          }
        }
      }
    }
  }
}