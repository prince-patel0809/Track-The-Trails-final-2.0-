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
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
</head>

<body style="margin:0; padding:0; background:#f4f6f8; font-family:Arial, sans-serif;">

<table width="100%" cellpadding="0" cellspacing="0">
<tr>
<td align="center">

<!-- CARD -->
<table width="500" style="background:#ffffff; margin-top:30px; border-radius:10px; padding:20px; box-shadow:0 2px 10px rgba(0,0,0,0.1);">

<tr>
<td align="center">
<h2 style="color:#2c3e50;">🚀 TRACK THE TRAILS</h2>
</td>
</tr>

<tr>
<td>
<p style="font-size:16px; color:#333;">Hello 👋,</p>
<p style="font-size:15px; color:#555;">
You have been assigned a new task.
</p>
</td>
</tr>

<tr>
<td style="background:#f8f9fa; padding:15px; border-radius:8px;">
<p><b>📁 Project:</b> `+projectName+`</p>
<p><b>👤 Assigned by:</b> `+ownerName+`</p>
</td>
</tr>

<tr><td height="15"></td></tr>

<tr>
<td style="background:#eef2ff; padding:15px; border-radius:8px;">
<p><b>📝 Title:</b> `+taskTitle+`</p>
<p><b>📄 Description:</b> `+description+`</p>
</td>
</tr>

<tr><td height="20"></td></tr>

<tr>
<td align="center">
<a href="#" style="
	background:#4f46e5;
	color:#ffffff;
	padding:12px 20px;
	text-decoration:none;
	border-radius:6px;
	font-size:14px;
">
View Task
</a>
</td>
</tr>

<tr><td height="20"></td></tr>

<tr>
<td align="center" style="font-size:12px; color:#999;">
© TRACK THE TRAILS — All rights reserved
</td>
</tr>

</table>

</td>
</tr>
</table>

</body>
</html>
`)

	d := gomail.NewDialer(
		host,
		587,
		from,
		password,
	)

	return d.DialAndSend(m)
}
