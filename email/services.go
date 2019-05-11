package email

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type DeliveryService struct {
	Host string
	Port int
	User string
	Pass string
}

// TODO: move this somewhere else?
func (es *DeliveryService) SendFeedbackEmail(emailAddress string, feedbackLink string) error {
	client := gomail.NewDialer(es.Host, es.Port, es.User, es.Pass)
	message := gomail.NewMessage()
	message.SetHeader("From", fmt.Sprintf("Anon Solicitor <%v>", es.User))
	message.SetHeader("To", emailAddress)
	body := fmt.Sprintf("<html><body><h3>Hey! We'd like to hear what you think!</h3><p>No worries - it's totally anonymous! Click <a href='%v'>here</a> to submit your feedback and see what everyone else thought!</p><p>Click <a href='%v'>here</a> to let us know that you didn't attend.</p><p>Thanks so much!</p></body></html>", feedbackLink, feedbackLink+"/absent")

	message.SetHeader("Title", fmt.Sprintf("You've been invited to give anonymous feedback about: %v", "test event"))
	message.SetBody("text/html", body)

	if err := client.DialAndSend(message); err != nil {
		log.Printf("failed to send email. Error: %v", err)

		return err
	}

	return nil
}
