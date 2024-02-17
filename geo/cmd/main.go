package main

import (
	"geo/internal/controller"
	"geo/internal/grpc/geo"
	"geo/internal/router"
	"geo/internal/service"
	pbgeo "geo/protos/gen/go"
	"github.com/go-chi/chi"
	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"sync"
)

func main() {

	con, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		log.Fatal(err)
	}

	queM, err := rabbitmq.NewRabbitMQ(con)
	if err != nil {
		log.Fatal(err)
	}

	err = rabbitmq.CreateExchange(con, "limitreq", "topic")
	if err != nil {
		log.Fatal(err)
	}

	ns := service.NewGeoService()
	nh := controller.NewHandGeo(ns)
	r := router.StartRout(nh, queM)
	w := sync.WaitGroup{}
	w.Add(2)
	go func(r *chi.Mux) {
		defer w.Done()
		log.Println("запуск сервера 8081")
		http.ListenAndServe(":8081", r)
	}(r)

	go func() {
		defer w.Done()
		listen, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Ошибка при прослушивании порта: %v", err)
		}

		server := grpc.NewServer()
		pbgeo.RegisterGeoServiceServer(server, &geo.ServerGeo{})

		log.Println("Запуск gRPC сервера geo...")
		if err := server.Serve(listen); err != nil {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()

	w.Wait()
}
