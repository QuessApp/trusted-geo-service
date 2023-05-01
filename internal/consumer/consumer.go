package consumer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/quessapp/toolkit/crypto"
	"github.com/quessapp/toolkit/entities"
	"github.com/quessapp/toolkit/i18n"
	"github.com/quessapp/toolkit/queue"
	"github.com/streadway/amqp"
)

// Consume Consumes from queue then publishes messages.
func Consume(geoServiceCh, emailServiceCh *amqp.Channel, client *sesv2.Client, geoServiceQueueName, emailServiceQueueName, cipherKey, mailFrom string) {
	msgs, err := queue.Consume(geoServiceCh, geoServiceQueueName)

	if err != nil {
		log.Fatalln(err)
	}

	for msg := range msgs {
		log.Printf("Received message %s \n", msg.Body)

		decryptedMessage, err := crypto.Decrypt(string(msg.Body), cipherKey)

		if err != nil {
			log.Fatalln(err)
			return
		}

		splittedMessage := strings.Split(decryptedMessage, "-")

		userEmail := splittedMessage[0]
		ip := splittedMessage[1]
		locale := splittedMessage[2]

		log.Printf("Searching for IP: %s \n", ip)

		r, err := http.Get("https://geolocation-db.com/json/" + ip)

		if err != nil {
			log.Println(err)
			return
		}

		b, err := io.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			return
		}

		location := entities.Location{}

		if err := json.Unmarshal(b, &location); err != nil {
			log.Printf("Error unmarshalling: %s \n", err)
			return
		}

		if err := msg.Ack(true); err != nil {
			log.Fatalln(err)
			return
		}

		log.Printf("Acked message %s and sending email to %s \n", msg.Body, userEmail)

		email := entities.Email{
			To:      userEmail,
			Subject: i18n.Translate(locale, "emails_unkown_login_attempt_subject"),
			Body:    fmt.Sprintf("%s %s, %s - %s", i18n.Translate(locale, "emails_unkown_login_attempt_body"), location.City, location.State, location.CountryName),
		}

		emailParsed, err := json.Marshal(email)

		if err != nil {
			log.Fatalf("fail to marshal %s", err)
			return
		}

		if err := queue.Publish(emailServiceCh, emailServiceQueueName, cipherKey, emailParsed); err != nil {
			log.Fatalf("fail to publish %s", err)
			return
		}
	}
}
