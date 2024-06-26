package notification

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
      :root {
        --font-family: "Inter var", sans-serif;
        --font-feature-settings: "cv02", "cv03", "cv04", "cv11";
        --primary-color: #f59e0b;
        --text-color: #4b5563;
      }

      .p-menubar {
        padding: 0.5rem;
        background: #f9fafb;
        color: #4b5563;
        border: 1px solid #e5e7eb;
        border-radius: 6px;
      }

      .h-2rem {
        height: 2rem !important;
      }

      .p-card {
        background: #ffffff;
        color: #4b5563;
        box-shadow: 0 2px 1px -1px rgba(0, 0, 0, 0.2), 0 1px 1px 0 rgba(0, 0, 0, 0.14), 0 1px 3px 0 rgba(0, 0, 0, 0.12);
        border-radius: 6px;
      }

      .p-component {
        font-family: var(--font-family);
        font-feature-settings: var(--font-feature-settings);
        font-size: 1rem;
        font-weight: normal;
      }

      img {
        background-image: url('/src/assets/temp.jpg');
        background-size: 518px 180px;
        background-repeat: no-repeat;
      }

      .p-card .p-card-body {
        padding: 1.25rem;
      }

      .p-card .p-card-title {
        font-size: 1.5rem;
        font-weight: 700;
        margin-bottom: 0.5rem;
      }

      .p-card .p-card-subtitle {
        font-weight: 400;
        margin-bottom: 0.5rem;
        color: #6b7280;
      }

      .pi {
        font-family: 'primeicons';
        font-style: normal;
        font-weight: normal;
        font-variant: normal;
        text-transform: none;
        line-height: 1;
        display: inline-block;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
      }

      .pi-map-marker:before {
        content: "📍";
      }

      .p-card .p-card-content {
        padding: 1.25rem 0;
      }

      .justify-content-around {
        justify-content: space-around !important;
      }

      .flex {
        display: flex !important;
      }

      .m-0 {
        margin: 0rem !important;
      }

      .mt-1 {
        margin-top: 0.25rem !important;
      }

      .flex-column {
        flex-direction: column !important;
      }

      .text-center {
        text-align: center !important;
      }

      .text-right {
        text-align: right !important;
      }

      .p-divider.p-divider-horizontal:before {
        border-top: 1px solid #e5e7eb;
      }

      .p-divider.p-divider-solid.p-divider-horizontal:before {
        border-top-style: solid;
      }

      .p-divider-horizontal:before {
        position: absolute;
        display: block;
        top: 50%;
        left: -1.25rem;
        width: 100%;
        content: "";
      }

      .p-divider.p-divider-horizontal {
        margin: 1.25rem 0;
        padding: 0 1.25rem;
      }

      .p-divider-horizontal {
        display: flex;
        width: 100%;
        position: relative;
        align-items: center;
      }

      .p-card .p-card-footer {
        padding: 1.25rem 0 0 0;
      }

      .gap-3 {
        gap: 1rem !important;
      }

      .w-full {
        width: 100% !important;
      }

      .p-button {
        color: #ffffff;
        background: #f59e0b;
        border: 1px solid #f59e0b;
        padding: 0.75rem 1.25rem;
        font-size: 1rem;
        transition: background-color 0.2s, color 0.2s, border-color 0.2s, box-shadow 0.2s;
        border-radius: 6px;
        outline-color: transparent;
        display: inline-flex;
        cursor: pointer;
        user-select: none;
        align-items: center;
        text-align: center;
        overflow: hidden;
        position: relative;
      }

      .p-button .p-button-label {
        transition-duration: 0.2s;
      }

      .p-button-label {
        flex: 1 1 auto;
        font-weight: 600;
      }

      .box-listing {
        width: 30rem;
      }

      .template {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
      }

      .header-logo {
        width: 100%;
        padding: 0.5rem;
        margin-bottom: 1rem;
        background: #f9fafb;
        color: #4b5563;
        border: 1px solid #e5e7eb;
        border-radius: 6px;
        justify-content: center;
        display: flex;
      }

      text {
        font-family: sans-serif;
      }

      .logo-text {
        display: flex;
        align-self: center;
        flex-direction: row;
      }

      .user-greeting {
        margin: 0 0 0 3rem;
        align-self: center;
        font-weight: 700;
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
	processingDuration := time.Since(listing.CreatedAt)
	client.TrackMetric("ListingCreationToNotificationDuration", processingDuration.Seconds())
	client.TrackMetric("UsersNotified", float64(users))
	usersNotifiedEvent := appinsights.NewEventTelemetry("UsersNotifiedEvent")
	usersNotifiedEvent.Properties["ListingID"] = listing.ID
	usersNotifiedEvent.Properties["UserCount"] = strconv.Itoa(users)
	usersNotifiedEvent.Properties["ProcessingDuration"] = processingDuration.String()
	client.Track(usersNotifiedEvent)
	client.Channel().Flush()
}
