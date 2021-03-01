package utils

import "testing"

func TestSendMail(t *testing.T) {
	a := NewAuthorizer("hculpan@pobox.com", "kQiV7Cuo")

	//The receiver needs to be in slice as the receive supports multiple receiver
	Receiver := []string{"harry@culpan.org"}

	Subject := "Testing another HTLML Email from golang"
	message := `
		<!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
		<html>
		<head>
		<meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1">
		</head>
		<body><img src="http://culpan.org/lordsofcrime/lordsofcrime.png">
		<p>This is a test email</p>
		</body>
		</html>
		`
	bodyMessage := a.WriteHTMLEmail(Receiver, "kingofcrime@culpan.org", Subject, message)

	a.SendMail(Receiver, "kingofcrime@culpan.org", Subject, bodyMessage)
}
