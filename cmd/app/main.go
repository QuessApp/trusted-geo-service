package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/quessapp/toolkit/queue"
	"github.com/quessapp/toolkit/ses"
	"github.com/quessapp/trusted-geo-service/internal/consumer"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file")
	}

	emailServiceQueueName := os.Getenv("EMAIL_SERVICE_QUEUE_NAME")
	emailServiceQueueURI := os.Getenv("EMAIL_SERVICE_QUEUE_URI")

	geoServiceQueueURI := os.Getenv("RABBITMQ_URI")
	geoServiceQueueName := os.Getenv("QUEUE_NAME")

	cipherKey := os.Getenv("CIPHER_KEY")

	mailFrom := os.Getenv("AWS_EMAIL_FROM")
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	region := os.Getenv("AWS_REGION")

	geoServiceConn, geoServiceCh := queue.Connect(geoServiceQueueURI)
	mailClient, err := ses.Configure(accessKey, secretKey, region)

	if err != nil {
		panic(err)
	}

	emailServiceConn, emailServiceCh := queue.Connect(emailServiceQueueURI)

	consumer.Consume(geoServiceCh, emailServiceCh, mailClient, geoServiceQueueName, emailServiceQueueName, cipherKey, mailFrom)

	defer geoServiceConn.Close()
	defer geoServiceCh.Close()

	defer emailServiceConn.Close()
	defer emailServiceCh.Close()
}
