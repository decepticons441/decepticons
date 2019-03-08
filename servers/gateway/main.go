package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"

	// "math/rand"
	"encoding/json"
	"sync/atomic"
	"time"

	"strings"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nehay100/assignments-nehay100/servers/gateway/handlers"
	"github.com/nehay100/assignments-nehay100/servers/gateway/indexes"
	"github.com/nehay100/assignments-nehay100/servers/gateway/models/users"
	"github.com/nehay100/assignments-nehay100/servers/gateway/sessions"
)

// func IndexHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello"))
// }

type Director func(r *http.Request)

func CustomDirector(targets []*url.URL, signingKey string, store sessions.Store) Director {
	var counter int32
	counter = 0
	state := &handlers.SessionState{}
	mx := sync.RWMutex{}
	mx.Lock()
	defer mx.Unlock()
	// log.Println(targets)
	return func(r *http.Request) {
		// _targets, _ := rc.Get("MessageAddresses").Result()
		// targets := strings.Split(_targets, ",")
		targ := targets[counter%int32(len(targets))]
		// log.Println(targ)

		atomic.AddInt32(&counter, 1) // note, to be extra safe, we’ll need to use mutexes
		counter++
		_, err := sessions.GetState(r, signingKey, store, state)
		if err != nil {
			r.Header.Del("X-User")
			fmt.Sprintf("Error getting session state/session unauthorized %v", err)
			return
		}
		// note the modulo (%) operator which maps some integer to range from 0 to
		// len(targets)
		j, err := json.Marshal(state.User)
		if err != nil {
			fmt.Sprintf("Error encoding session state user %v", err)
			return
		}
		r.URL.Host = targ.Host
		r.Host = targ.Host
		r.URL.Scheme = "http"
		r.Header.Add("X-User", string(j))
	}
}

//main is the main entry point for the server
func main() {
	// Read the ADDR environment variable to get the address
	// the server should listen on. If empty, default to ":80"
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	if len(tlsKeyPath) == 0 {
		log.Println("Can't Obtain TLSKEY environment variable")
		os.Exit(1)
	}

	if len(tlsCertPath) == 0 {
		log.Println("Can't Obtain TLSCERT environment variable")
		os.Exit(1)
	}
	sessionkey := os.Getenv("SESSIONKEY")
	if len(sessionkey) == 0 {
		sessionkey = "default"
	}
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "127.0.0.1:6379"
	}

	db, err := sql.Open("mysql", os.Getenv("DSN"))

	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	//for now, just ping the server to ensure we have
	//a live connection to it
	if err := db.Ping(); err != nil {
		fmt.Printf("error pinging database: %v\n", err)
	} else {
		fmt.Printf("successfully connected to db!\n")
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Printf("Pong: ", pong)
		fmt.Printf("error pinging redis: %v\n", err)
	} else {
		fmt.Printf("successfully connected to redis!\n")
	}
	redisStore := sessions.NewRedisStore(client, time.Hour)

	store := users.NewSqlStore(db)
	sh := &handlers.SessionHandler{
		SigningKey: sessionkey,
		Store:      redisStore,
		Users:      store,
		Trie:       *indexes.NewTrie(),
	}
	_, err = sh.Users.AddAllToTrie()
	if err != nil {
		fmt.Sprintf("User-Handler: Error when inserting new user to trie %v", err)
	}

	messageAddresses := strings.Split(os.Getenv("MESSAGE_ADDR"), ",")
	summaryAddresses := strings.Split(os.Getenv("SUMMARY_ADDR"), ",")
	var messageURLs []*url.URL
	for _, address := range messageAddresses {
		messageURL, err := url.Parse(address)
		if err != nil {
			fmt.Sprintf("Error when parsing message addresss into url format %v", err)
		}
		messageURLs = append(messageURLs, messageURL)
	}

	var summaryURLs []*url.URL
	for _, address := range summaryAddresses {
		summaryURL, err := url.Parse(address)
		if err != nil {
			fmt.Sprintf("Error when parsing summary addresss into url format %v", err)
		}
		summaryURLs = append(summaryURLs, summaryURL)
	}

	messageProxy := &httputil.ReverseProxy{Director: CustomDirector(messageURLs, sh.SigningKey, sh.Store)}
	summaryProxy := &httputil.ReverseProxy{Director: CustomDirector(summaryURLs, sh.SigningKey, sh.Store)}

	// Create a new mux for the web server.
	// muxHandler := http.NewServeMux()

	r := mux.NewRouter()
	// r.HandleFunc("/products/{key}", ProductHandler)
	// r.HandleFunc("/articles/{category}/", ArticlesCategoryHandler)
	// r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)

	// r.HandleFunc("/v1/users", sh.UsersHandler)
	// r.HandleFunc("/v1/users/", sh.SpecificUserHandler)
	// r.HandleFunc("/v1/sessions", sh.SessionsHandler)
	// r.HandleFunc("/v1/sessions/", sh.SpecificSessionHandler)

	// r.Handle("/v1/channels", messageProxy) // register the proxies
	// r.Handle("/v1/channels/", messageProxy)
	// r.Handle("/v1/channels/members", messageProxy)
	// r.Handle("/v1/messages/", messageProxy)

	// r.Handle("/v1/summary", summaryProxy)

	// Tell the mux to call your handlers.SummaryHandler function
	// when the "/v1/summary" URL path is requested.

	ws := handlers.SocketStore{
		Connections: make(map[int64]*websocket.Conn),
		Sh:          sh,
	}

	rabbit := os.Getenv("RABBIT")
	conn, err := amqp.Dial("amqp://rabbit:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbit, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msg, err := ch.Consume(
		q.Name,
		"",    // Consumer
		true,  // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)

	go ws.SendMessages(msg)
	r.HandleFunc("/v1/ws", ws.WebSocketConnectionHandler)

	r.HandleFunc("/v1/users", sh.UsersHandler)
	r.HandleFunc("/v1/users/{id}", sh.SpecificUserHandler)
	r.HandleFunc("/v1/sessions", sh.SessionsHandler)
	r.HandleFunc("/v1/sessions/{id}", sh.SpecificSessionHandler)

	r.Handle("/v1/channels", messageProxy) // register the proxies
	r.Handle("/v1/channels/{id}", messageProxy)
	r.Handle("/v1/channels/{channelID}/members", messageProxy)
	r.Handle("/v1/messages/{id}", messageProxy)

	r.Handle("/v1/summary", summaryProxy)
	wrappedMux := handlers.NewCorsHandler(r)

	// Start a web server listening on the address you read from
	// the environment variable, using the mux you created as
	// the root handler. Use log.Fatal() to report any errors
	// that occur when trying to start the web server.
	fmt.Printf("server is listening on %s...\n", addr)
	// log.Fatal(http.ListenAndServe(addr, mux))

	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
