package notification

import (
	"bytes"
	"fmt"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"html/template"
	"strconv"
	"time"
)

const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            border: 1px solid #ddd;
            border-radius: 8px;
        }
        .header {
            text-align: center;
            margin-bottom: 20px;
        }
        .header img {
            max-width: 100%;
            border-radius: 8px;
        }
        .content {
            margin-bottom: 20px;
        }
        .content h2 {
            margin-top: 0;
        }
        .content p {
            margin: 5px 0;
        }
        .button {
            text-align: center;
        }
        .button a {
            background-color: #28a745;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <img src="{{.PreviewImageURL}}" alt="Listing Image">
        </div>
        <div class="content">
            <h2>{{.Title}}</h2>
            <p><strong>Location:</strong> {{.Address}}, {{.City}}, {{.Country}}, {{.PostalCode}}</p>
            <p><strong>Room Count:</strong> {{.RoomCount}}</p>
            <p><strong>Size:</strong> {{.SquareMeters}} m²</p>
            <p><strong>Availability Date:</strong> {{.AvailabilityDate}}</p>
            <p><strong>Year Built:</strong> {{.YearBuilt}}</p>
            <p><strong>HWG Energy Class:</strong> {{.HwgEnergyClass}}</p>
            <p><strong>FGEE Energy Class:</strong> {{.FgeeEnergyClass}}</p>
            <p><strong>Listing Type:</strong> {{.ListingType}}</p>
            {{if .RentPricePerMonth}}<p><strong>Rent Price Per Month:</strong> €{{.RentPricePerMonth}}</p>{{end}}
            {{if .CooperativeShare}}<p><strong>Cooperative Share:</strong> €{{.CooperativeShare}}</p>{{end}}
            {{if .SalePrice}}<p><strong>Sale Price:</strong> €{{.SalePrice}}</p>{{end}}
            {{if .AdditionalFees}}<p><strong>Additional Fees:</strong> €{{.AdditionalFees}}</p>{{end}}
        </div>
        <div class="button">
            <a href="{{.DetailsURL}}" target="_blank">View Details</a>
        </div>
    </div>
</body>
</html>
`

func generateEmailContent(listing models.Listing) (string, error) {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, listing); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SendNotification(listing models.Listing, emails []string) error {
	content, err := generateEmailContent(listing)
	if err != nil {
		return fmt.Errorf("failed to generate email content: %w", err)
	}

	subject := fmt.Sprintf("New GeWo in %s matches your preferences!", listing.City)

	return sendHtmlEmail(emails, subject, content)
}

func PublishTelemetry(client appinsights.TelemetryClient, listing *models.Listing, users int) {
	processingDuration := time.Now().Sub(listing.CreatedAt)
	client.TrackMetric("ListingCreationToNotificationDuration", processingDuration.Seconds())
	client.TrackMetric("UsersNotified", float64(users))
	usersNotifiedEvent := appinsights.NewEventTelemetry("UsersNotifiedEvent")
	usersNotifiedEvent.Properties["ListingID"] = listing.ID
	usersNotifiedEvent.Properties["UserCount"] = strconv.Itoa(users)
	usersNotifiedEvent.Properties["ProcessingDuration"] = processingDuration.String()
	client.Track(usersNotifiedEvent)
	client.Channel().Flush()
}
