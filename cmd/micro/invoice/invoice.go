package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const verison = "1..0.0"

type config struct {
	port int
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	frontend string
}
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
}

func (app *application) Serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.invoceRoute(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting Back end Server In %d", app.config.port)
	return srv.ListenAndServe()
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 5000, "Server Port To Listen On")
	flag.StringVar(&cfg.smtp.host, "smtphost", "sandbox.smtp.mailtrap.io", "smpt host")
	flag.IntVar(&cfg.smtp.port, "smtpport", 587, "smpt port")
	flag.StringVar(&cfg.smtp.username, "smtpusername", "2e2238e526e2c7", "smpt username")
	flag.StringVar(&cfg.smtp.password, "smtppassword", "4595f6612f03b7", "smpt password")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "domain frontend")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorfoLog,
		version:  verison,
	}

	app.CreateDirIfNotExist("./invoices")
	err := app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
