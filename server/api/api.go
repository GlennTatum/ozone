package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"ozone/util"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	gocql "github.com/gocql/gocql"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const REDIS_URL = "127.0.0.1:6379"
const CASSANDRA_URL = "127.0.0.1:9042"

type App struct {
	r       *mux.Router
	Session *scs.SessionManager
	K       *util.ResourceManagerClient
	Db      *gocql.Session
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

func NewApp() (*App, error) {

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

	var mgmt *util.ResourceManagerClient
	_, ok := os.LookupEnv("PRODUCTION")
	if ok {
		m, err := util.NewInClusterResourceManagerClient()
		mgmt = m
		if err != nil {
			return nil, err
		}
	} else {
		kconfig, ok := os.LookupEnv("KUBECONFIG")
		if !ok {
			return nil, errors.New("No Kubeconfig found")
		}
		m, err := util.NewOutOfClusterResourceManagerClient(kconfig)
		mgmt = m
		if err != nil {
			return nil, err
		}
	}

	return &App{
		r:       r,
		K:       mgmt,
		Session: sessionManager,
		Db:      db,
		logger:  sugar,
		cache:   pool,
	}, nil
}

func Exec() {
	app, err := NewApp()

	app.r.Use(app.Session.LoadAndSave)
	app.r.HandleFunc("/", app.HomeRoute)

	e := app.r.PathPrefix("/events").Subrouter()
	e.Use(app.Authz)
	e.HandleFunc("/e/{id}", app.EventRoute) // .Methods("POST") setup cookieJar in postman

	app.logger.Debugw("starting server...")
	srv := &http.Server{
		Handler:      app.r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if err != nil {
		panic(err)
	}
	log.Fatal(srv.ListenAndServe())
}
