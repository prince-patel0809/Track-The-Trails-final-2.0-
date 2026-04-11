package utils

import (
	"log"
	"strconv"
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
	portStr := config.GetEnv("SMTP_PORT")

	// ===== DEBUG ENV =====
	log.Println("SMTP HOST:", host)
	log.Println("SMTP EMAIL:", from)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 465
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "TRACK THE TRAILS - New Task Assigned 🚀")

	m.SetBody("text/html", `<!DOCTYPE html>
<html>
<body style="font-family:Arial; background:#f4f6f8; padding:20px;">
<div style="max-width:500px; margin:auto; background:white; padding:20px; border-radius:10px;">
<h2>🚀 TRACK THE TRAILS</h2>

<p>Hello 👋,</p>
<p>You have been assigned a new task.</p>

<hr>

<p><b>📁 Project:</b> `+projectName+`</p>
<p><b>👤 Assigned by:</b> `+ownerName+`</p>

<hr>

<p><b>📝 Title:</b> `+taskTitle+`</p>
<p><b>📄 Description:</b> `+description+`</p>

<hr>

<p style="color:gray;">Login to TRACK THE TRAILS to manage your task.</p>

</div>
</body>
</html>`)

	d := gomail.NewDialer(host, port, from, password)

	// 🔥 IMPORTANT FIX FOR RENDER
	d.SSL = false

	// ===== SEND EMAIL =====
	err = d.DialAndSend(m)
	if err != nil {
		log.Println("❌ EMAIL ERROR:", err)
		return err
	}

	log.Println("✅ Email sent successfully to:", toEmail)
	return nil
}
