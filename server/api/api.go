package api

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	gocql "github.com/gocql/gocql"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	gocqlx "github.com/scylladb/gocqlx/v2"
	"go.uber.org/zap"
)

const REDIS_URL = "127.0.0.1:6379"
const CASSANDRA_URL = "127.0.0.1:9042"

type App struct {
	r       *mux.Router
	session *scs.SessionManager
	db      *gocqlx.Session
	logger  *zap.SugaredLogger
	cache   *redis.Pool
}

func cassandraSession(cluster *gocql.ClusterConfig) *gocqlx.Session {
	conn, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal(err)
	}
	return &conn
}

func NewApp() *App {

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", REDIS_URL)
		},
	}

	cluster := gocql.NewCluster(CASSANDRA_URL)
	cluster.Keyspace = "main"
	db := cassandraSession(cluster)

	sessionManager := scs.New()
	sessionManager.Store = redisstore.New(pool)

	r := mux.NewRouter()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	return &App{
		r:       r,
		session: sessionManager,
		db:      db,
		logger:  sugar,
		cache:   pool,
	}
}

func Exec() {
	app := NewApp()

	app.r.Use(app.session.LoadAndSave)
	app.r.Use(app.Authz)

	app.r.HandleFunc("/healthz", app.healthz)

	app.logger.Debugw("starting server...")
	srv := &http.Server{
		Handler:      app.r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
