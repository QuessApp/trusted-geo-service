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

		log.Printf("Acked message %s \n", msg.Body)

		log.Printf("Sending email to %s \n", userEmail)

		email := entities.Email{
			To:      userEmail,
			Subject: "Novo login desconhecido",
			Body:    fmt.Sprintf("Percebemos que você fez login em um novo dispositivo. Oriundo de %s, %s no %s. Se foi você, não há nada para você fazer agora. Se não foi você, entre em contato com a gente imediatamente.", location.City, location.State, location.CountryName),
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
