package main

import (
	"fmt"
	"regexp"

	"github.com/resendlabs/resend-go"
)

type EmailService struct {
	DbConnectionString string
}

type Email struct {
	Sender    string
	Recipient string
	Message   string
}

func (service *EmailService) Send(email Email, protocol string) {

	if !service.Validate(email) {
		fmt.Println("Invalid email")
		return
	}

	switch protocol {
	case "SMTP":
		fmt.Println("Sending email via SMTP... (not really)")

	case "IMAP":
		fmt.Println("Sending email via IMAP... (not really)")

	case "POP3":
		fmt.Println("Sending email via POP3... (not really)")

	case "RESEND":
		// contact me for api key, it works :-) (slam23@vse.cz)
		apiKey := "re_123"
		client := resend.NewClient(apiKey)

		// Send
		params := &resend.SendEmailRequest{
			From:    email.Sender,
			To:      []string{email.Recipient},
			Subject: "Hello world",
			Html:    "<strong>" + email.Message + "</strong>",
		}
		sent, err := client.Emails.Send(params)
		if err != nil {
			panic(err)
		}

		fmt.Println("Sending email via resend.com - MsgId: " + sent.Id)

	default:
		fmt.Println("Sending email via default protocol (SMTP)... (not really)")
	}

	service.Store(email)
}

// Validate email address of Sender and Recipient
func (service *EmailService) Validate(emailToValidate Email) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(emailToValidate.Sender) && emailRegex.MatchString(emailToValidate.Recipient)
}

// Store an email (dummy)
func (service *EmailService) Store(email Email) {
	fmt.Printf("Storing email... (not really)\n")
	fmt.Printf("\tFrom: %s\n\tTo: %s\n\tMessage: %s\n ", email.Sender, email.Recipient, email.Message)
}

func main() {
	// Initialize EmailService with a database connection string
	emailService := EmailService{
		DbConnectionString: "your_db_connection_string",
	}

	// Initialize an Email
	email := Email{
		Sender:    "no-reply@tapeer.cz",
		Recipient: "emptykiwi@gmail.com",
		Message:   "Hello, this is a test email!",
	}

	// Send an email using the EmailService
	emailService.Send(email, "RESEND")
}
