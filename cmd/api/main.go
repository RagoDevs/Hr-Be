package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/Hopertz/Hr-Be/config"
	db "github.com/Hopertz/Hr-Be/internal/db/sqlc"
)

type application struct {
	config config.Config
	wg     sync.WaitGroup
	store  db.Store
}

const version = "1.0.0"

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		petstore.swagger.io
//	@BasePath	/v2

func main() {
	var cfg config.Config

	flag.IntVar(&cfg.PORT, "port", 8668, "API server port")
	flag.StringVar(&cfg.ENV, "env", "development", "Environment (development|Staging|production")
	flag.StringVar(&cfg.DB.DSN, "db-dsn", os.Getenv("HR_DB_DSN"), "PostgresDB DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max ilde connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection  connections")

	flag.Parse()

	conn, err := config.OpenDB(cfg)
	if err != nil {
		log.Fatal(err, nil)
	}

	defer conn.Close()

	slog.Info("database connection established")

	app := &application{
		config: cfg,
		store:  db.NewStore(conn),
	}

	slog.Info("Starting server on", "port", cfg.PORT)
	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}

}
