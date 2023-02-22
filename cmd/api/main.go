package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-playground/form"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"log"
	"net/http"
	"os"
	"s_renovation.net/internal/data"
	"time"
)

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	models        data.Models
	formDecoder   *form.Decoder
	templateCache map[string]*template.Template
}

type config struct {
	port int
	env  string
	db   struct {
		dsn         string
		maxOpenConn uint64
		maxIdleTime time.Duration
	}
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production")

	flag.StringVar(&cfg.db.dsn, "mongo_uri", os.Getenv("MONGO_DB_DSN"), "db connection string")
	flag.Uint64Var(&cfg.db.maxOpenConn, "maxOpenConn", uint64(100), "maximum number of open connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "maxIdleTime", time.Duration(10), "maximum idle time of one connection")
	flag.Parse()

	infoLog := log.New(os.Stdout, "info\t", log.LstdFlags)
	errorLog := log.New(os.Stdout, "error\t", log.LstdFlags|log.Lshortfile)

	client, err := openDB(cfg)
	if err != nil {
		errorLog.Fatalf("Could not connect to MongoDB Atlas: %v", err)
	}
	infoLog.Println("Connected to Mongo Atlas!")
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			errorLog.Fatalf("Couldn't close the database connection, due to: %v", err)
		}
	}()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	
	app := application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		models:        data.NewModels(client),
		formDecoder:   formDecoder,
		templateCache: templateCache,
	}
	srv := http.Server{
		Addr:     fmt.Sprintf("localhost:%v", cfg.port),
		Handler:  app.Router(),
		ErrorLog: errorLog,
	}
	infoLog.Printf("staring %v server on %v", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatalln(err)

}

func openDB(cfg config) (*mongo.Client, error) {
	ClientOptions := options.ClientOptions{
		MaxPoolSize:     &cfg.db.maxOpenConn,
		MaxConnIdleTime: &cfg.db.maxIdleTime,
	}

	client, err := mongo.Connect(context.TODO(), ClientOptions.ApplyURI(cfg.db.dsn))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

///написать валидатор бизнес логики
//написать парсинг темплейтов
