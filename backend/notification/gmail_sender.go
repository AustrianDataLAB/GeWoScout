package notification

import (
	"bytes"
	"fmt"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strings"
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

func GenerateEmailContent(listing models.Listing) (string, error) {
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

func GmailDemo() {
	rent := 1200
	listing := models.Listing{
		ID:                 "1",
		Title:              "Beautiful Apartment",
		HousingCooperative: "Example Cooperative",
		ProjectID:          "123",
		ListingID:          "456",
		Country:            "Wonderland",
		City:               "Fictitious City",
		PostalCode:         "12345",
		Address:            "123 Example Street",
		RoomCount:          3,
		SquareMeters:       75,
		AvailabilityDate:   "2024-07-01",
		YearBuilt:          2000,
		HwgEnergyClass:     "A+",
		FgeeEnergyClass:    "A",
		ListingType:        "rent",
		RentPricePerMonth:  &rent,
		DetailsURL:         "http://example.com/details",
		PreviewImageURL:    "http://example.com/image.jpg",
	}

	content, err := GenerateEmailContent(listing)
	if err != nil {
		log.Fatal(err)
	}

	sendEmail(content)
}

func sendEmail(content string) {
	_ = os.Setenv("GMAIL_USER", "noreply.gewoscout@gmail.com")
	//_ = os.Setenv("GMAIL_PASS", "<FILL IN HERE>")

	// Get credentials from environment variables.
	email := os.Getenv("GMAIL_USER")
	password := os.Getenv("GMAIL_PASS")

	// Set up authentication information.
	auth := smtp.PlainAuth("", email, password, "smtp.gmail.com")

	// Set up email details.
	from := email
	to := []string{"asdf@gmail.com"}
	headers := "MIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\""
	subject := "HTML test 2"
	body := content

	// Concatenate email parts.
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s\r\n\r\n%s", strings.Join(to, ", "), subject, headers, body))

	// Set up the SMTP server configuration.
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		to,
		msg,
	)

	// Check for errors.
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
}
