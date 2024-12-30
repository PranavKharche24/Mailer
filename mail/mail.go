package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/pranavKharche24/mail/utils"
)

func SendMailSimpleHTML(to []string, subject, htmlFile string, cc, bcc []string, attachments []string) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(htmlFile)
	if err != nil {
		return fmt.Errorf("error parsing HTML file: %v", err)
	}
	if err := t.Execute(&body, struct{ Name string }{Name: "Pranav"}); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}
	return sendEmailWithAttachments(to, subject, body.String(), cc, bcc, attachments, true)
}

func SendMailPlain(to []string, subject, msg string, cc, bcc []string, attachments []string) error {
	return sendEmailWithAttachments(to, subject, msg, cc, bcc, attachments, false)
}

func sendEmailWithAttachments(to []string, subject, content string, cc, bcc []string, attachments []string, isHTML bool) error {
	auth := smtp.PlainAuth("", utils.FromEmail, utils.FromPass, "smtp.gmail.com")

	recipients := append(append([]string{}, to...), cc...)
	recipients = append(recipients, bcc...)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fmt.Fprintf(&buf, "From: %s\r\n", utils.FromEmail)
	fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(to, ","))
	if len(cc) > 0 {
		fmt.Fprintf(&buf, "Cc: %s\r\n", strings.Join(cc, ","))
	}
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, "Content-Type: multipart/mixed; boundary=%s\r\n\r\n", writer.Boundary())

	bodyHeader := make(textproto.MIMEHeader)
	if isHTML {
		bodyHeader.Set("Content-Type", `text/html; charset="UTF-8"`)
	} else {
		bodyHeader.Set("Content-Type", `text/plain; charset="UTF-8"`)
	}
	bodyPart, err := writer.CreatePart(bodyHeader)
	if err != nil {
		return fmt.Errorf("error creating body part: %v", err)
	}
	if _, err := bodyPart.Write([]byte(content)); err != nil {
		return fmt.Errorf("error writing body part: %v", err)
	}

	for _, filePath := range attachments {
		if filePath == "" {
			continue
		}
		if err := utils.AttachFile(writer, filePath); err != nil {
			return fmt.Errorf("error attaching file: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	if err := smtp.SendMail("smtp.gmail.com:587", auth, utils.FromEmail, recipients, buf.Bytes()); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}
