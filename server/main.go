package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const REDIS_URL = "127.0.0.1:6379"
const CASSANDRA_URL = "127.0.0.1:9042"

type App struct {
	r       *mux.Router
	session *scs.SessionManager
	db      *gocql.Session
	logger  *zap.SugaredLogger
	cache   *redis.Pool
}

func cassandraSession(cluster *gocql.ClusterConfig) *gocql.Session {
	conn, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	return conn
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
	r.Use(sessionManager.LoadAndSave)

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

func main() {
	app := NewApp()
	app.r.HandleFunc("/put", app.putHandler)
	app.r.HandleFunc("/get", app.getHandler)

	app.logger.Debugw("starting server...")
	srv := &http.Server{
		Handler:      app.r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (app *App) putHandler(w http.ResponseWriter, r *http.Request) {
	app.session.Put(r.Context(), "message", "Hello from a session!")
}

func (app *App) getHandler(w http.ResponseWriter, r *http.Request) {
	msg := app.session.GetString(r.Context(), "message")
	io.WriteString(w, msg)
}
