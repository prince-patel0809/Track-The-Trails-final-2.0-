package utils

import (
	"track-the-trails/config"

	"gopkg.in/gomail.v2"
)

func SendTaskAssignmentEmail(
	toEmail string,
	projectName string,
	ownerName string,
	taskTitle string,
	description string,
) error {

	from := config.GetEnv("SMTP_EMAIL")
	password := config.GetEnv("SMTP_PASSWORD")
	host := config.GetEnv("SMTP_HOST")

	m := gomail.NewMessage()

	// ===== EMAIL HEADER =====
	m.SetHeader("From", from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "TRACK THE TRAILS - New Task Assigned 🚀")

	// ===== EMAIL BODY =====
	m.SetBody("text/html", `
		<div style="font-family: Arial, sans-serif; padding: 20px;">
			
			<h1 style="color:#2c3e50;">🚀 TRACK THE TRAILS</h1>

			<h3>You have been assigned a new task!</h3>

			<hr/>

			<p><b>📁 Project:</b> `+projectName+`</p>
			<p><b>👤 Assigned by:</b> `+ownerName+`</p>

			<hr/>

			<h3>📝 Task Details</h3>
			<p><b>Title:</b> `+taskTitle+`</p>
			<p><b>Description:</b> `+description+`</p>

			<hr/>

			<p style="color:gray;">Please login to TRACK THE TRAILS to view and manage your task.</p>

		</div>
	`)

	d := gomail.NewDialer(
		host,
		587,
		from,
		password,
	)

	return d.DialAndSend(m)
}
