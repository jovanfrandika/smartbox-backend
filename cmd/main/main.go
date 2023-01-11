package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
	userContainer "github.com/jovanfrandika/smartbox-backend/pkg/user"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "test"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	config.Init()

	client, err := mongo.NewClient(options.Client().ApplyURI(config.Cfg.MongoDBUri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbName)
	r := chi.NewRouter()

	userRouter := chi.NewRouter()
	userContainer.Init(db, userRouter, config.Cfg)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	r.Mount("/user", userRouter)

	http.ListenAndServe(":8000", r)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
