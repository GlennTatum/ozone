package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"ozone/util"
	"strings"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	gocql "github.com/gocql/gocql"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type App struct {
	r       *mux.Router
	Session *scs.SessionManager
	K       *util.ResourceManagerClient
	Db      *gocql.Session
	logger  *zap.SugaredLogger
	cache   *redis.Pool
}

func cassandraSession(cluster *gocql.ClusterConfig) *gocql.Session {

	v, _ := os.LookupEnv("APP_ENV")
	if v == "PRODUCTION" {
		cluster_username, _ := os.LookupEnv("CASSANDRA_USERNAME")
		cluster_password, _ := os.LookupEnv("CASSANDA_PASSWORD")
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cluster_username, // Replace with your Cassandra username
			Password: cluster_password, // Replace with your Cassandra password
		}
	}
	conn, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func NewApp() (*App, error) {

	REDIS_URL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		return nil, errors.New("REDIS_URL not found")
	}
	CASSANDRA_URL, ok := os.LookupEnv("CASSANDRA_URL")
	if !ok {
		return nil, errors.New("CASSANDRA_URL not found")
	}
	hosts := strings.Split(CASSANDRA_URL, ",")

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", REDIS_URL)
		},
	}

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "main"
	db := cassandraSession(cluster)

	sessionManager := scs.New()
	sessionManager.Store = redisstore.New(pool)

	r := mux.NewRouter()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	var mgmt *util.ResourceManagerClient
	v, ok := os.LookupEnv("APP_ENV")
	if ok && (v == "PRODUCTION") {
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
	if err != nil {
		panic(err)
	}

	// Setup CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://code.gdg-rit.dev"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		Debug:            true,
	})

	app.r.Use(app.Session.LoadAndSave)
	app.r.HandleFunc("/", app.HomeRoute)

	e := app.r.PathPrefix("/events").Subrouter()
	e.Use(app.Authz)
	e.HandleFunc("/e/{id}", app.EventRoute).Methods("POST")

	app.logger.Debugw("starting server...")

	// Wrap router with CORS
	srv := &http.Server{
		Handler:      c.Handler(app.r),
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
