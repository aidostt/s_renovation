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
	"s_renovation.net/internal/jsonlog"
	"time"
)

type application struct {
	config        config
	logger        *jsonlog.Logger
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

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	client, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	logger.PrintInfo("database connection pool established", nil)
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logger.PrintFatal(err, nil) //nigga
		}
	}()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	formDecoder := form.NewDecoder()

	app := application{
		config: cfg,
		logger: logger,
		//infoLog:       infoLog, //nigga
		//errorLog:      errorLog, //nigga
		models:        data.NewModels(client),
		formDecoder:   formDecoder,
		templateCache: templateCache,
	}
	srv := http.Server{
		Addr:     fmt.Sprintf("localhost:%v", cfg.port),
		Handler:  app.Router(),
		ErrorLog: log.New(logger, "", 0), //nigga
	}

	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})

	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)

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
