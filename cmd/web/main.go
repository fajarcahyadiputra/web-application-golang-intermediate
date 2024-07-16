package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/fajarcahyadiputra/udemy-web-application/internal/driver"
	"github.com/fajarcahyadiputra/udemy-web-application/internal/models"
)

const verison = "1..0.0"
const cssVersion = "1"

var session *scs.SessionManager

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
	secrectkey string
	frontend   string
}
type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	DB            models.DBModel
	Session       *scs.SessionManager
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

	app.infoLog.Printf("Starting HTTP Server In %s Mode On Port %d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}

func main() {
	gob.Register(TransactionData{})
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Server Port To Listen On")
	flag.StringVar(&cfg.env, "env", "development", "Application Environment {development|prodyction|testing}")
	flag.StringVar(&cfg.db.dsn, "dsn", "root:12345678@tcp(localhost:3306)/learning_widgets?parseTime=true&tls=false", "DSN")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to API")
	flag.StringVar(&cfg.secrectkey, "secrectkey", "jdu73tdjruplcjry36ahsyebncmxkipe", "secrect key")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "domain frontend")

	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorfoLog.Fatal(err)
	}
	defer conn.Close()

	//set up session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = mysqlstore.New(conn)

	tc := make(map[string]*template.Template)
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorfoLog,
		templateCache: tc,
		version:       verison,
		DB: models.DBModel{
			DB: conn,
		},
		Session: session,
	}

	go app.ListenToWSChannel()

	err = app.Serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}

}
