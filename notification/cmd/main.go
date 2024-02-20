package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/kafkamq"
	"gitlab.com/ptflp/gopubsub/queue"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"log"
	"notification/internal/models"
	"notification/internal/service"
	"os"
)

func main() {

	var queM queue.MessageQueuer
	protocol := os.Getenv("BROKER")
	switch protocol {
	case "rabbit":
		log.Println("conn to rabbit")
		conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
		if err != nil {
			log.Fatal(err)
		}
		queR, err := rabbitmq.NewRabbitMQ(conn)
		if err != nil {
			log.Fatal(err)
		}

		err = rabbitmq.CreateExchange(conn, "limitreq", "topic")
		if err != nil {
			log.Fatal(err)
		}
		queM = queR
	case "kafka":
		log.Println("conn to kafka")
		queK, err := kafkamq.NewKafkaMQ("kafka:9092", "mG")
		if err != nil {
			log.Fatal(err)
		}

		queM = queK
	}

	msgs, err := queM.Subscribe("limitreq")
	if err != nil {
		log.Fatal(err)
	}

	d := &models.Data{}
	ns := service.NewNotifService()

	for msg := range msgs {

		err = queM.Ack(&msg)
		if err != nil {
			log.Println("not data msg")
		}
		err = json.Unmarshal(msg.Data, &d)
		if err != nil {
			log.Println("err unmarsh")
		}

		respM, _ := ns.SendItToTheEmail(d.Email)
		log.Println(respM)

		respPh, _ := ns.SentItToThePhone(d.Phone)
		log.Println(respPh)
	}

}
