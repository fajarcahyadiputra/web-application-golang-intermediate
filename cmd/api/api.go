package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fajarcahyadiputra/udemy-web-application/internal/driver"
	"github.com/fajarcahyadiputra/udemy-web-application/internal/models"
)

const verison = "1..0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	secrectkey string
	frontend   string
}
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	DB       models.DBModel
}

func (app *application) Serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting Back end Server In %s Mode On Port %d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4001, "Server Port To Listen On")
	flag.StringVar(&cfg.env, "env", "development", "Application Environment {development|prodyction|testing}")
	flag.StringVar(&cfg.db.dsn, "dsn", "root:12345678@tcp(localhost:3306)/learning_widgets?parseTime=true&tls=false", "DSN")
	flag.StringVar(&cfg.smtp.host, "smtphost", "sandbox.smtp.mailtrap.io", "smpt host")
	flag.IntVar(&cfg.smtp.port, "smtpport", 587, "smpt port")
	flag.StringVar(&cfg.smtp.username, "smtpusername", "2e2238e526e2c7", "smpt username")
	flag.StringVar(&cfg.smtp.password, "smtppassword", "4595f6612f03b7", "smpt password")
	flag.StringVar(&cfg.secrectkey, "secrectkey", "jdu73tdjruplcjry36ahsyebncmxkipe", "secrect key")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "domain frontend")

	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	con, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorfoLog.Fatal(err)
	}
	defer con.Close()
	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorfoLog,
		version:  verison,
		DB: models.DBModel{
			DB: con,
		},
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
