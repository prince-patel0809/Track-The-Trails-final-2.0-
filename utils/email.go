package utils

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendTaskAssignmentEmail(toEmail, projectName, ownerName, taskTitle, description string) error {

	from := mail.NewEmail("TRACK THE TRAILS", "propatel0809@gmail.com")
	to := mail.NewEmail("User", toEmail)

	subject := "TRACK THE TRAILS - New Task Assigned 🚀"

	content := `
🚀 TRACK THE TRAILS

Project: ` + projectName + `
Assigned by: ` + ownerName + `

Task: ` + taskTitle + `
Description: ` + description + `
`

	message := mail.NewSingleEmail(from, subject, to, content, content)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)
	if err != nil {
		log.Println("SENDGRID ERROR:", err)
		return err
	}

	log.Println("Email sent:", response.StatusCode)
	return nil
}
