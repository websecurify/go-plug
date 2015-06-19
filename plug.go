package plug

// ---
// ---
// ---

import (
	"os"
	"log"
	"net/http"
	
	// ---
	
	"gopkg.in/mgo.v2"
	
	// ---
	
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/gorilla/handlers"
)

// ---
// ---
// ---

var G struct {
	C *mgo.Collection
}

// ---
// ---
// ---

func Getenv(name string, defaultValue string) (string) {
	value := os.Getenv(name)
	
	// ---
	
	if value != "" {
		return value
	} else {
		return defaultValue
	}
}

// ---
// ---
// ---

func Info(v ...interface{}) {
	log.Println(v...)
}

func Error(v ...interface{}) {
	log.Println(v...)
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// ---
// ---
// ---

func Run(receiver interface{}, name string) {
	s := rpc.NewServer()
	
	// ---
	
	s.RegisterCodec(json.NewCodec(), "application/json")
	
	// ---
	
	s.RegisterService(receiver, name)
	
	// ---
	
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, s))
	
	// ---
	
	address := Getenv("ADDRESS", ":8080")
	
	// ---
	
	Info("starting server at", address)
	
	// ---
	
	Fatal(http.ListenAndServe(address, nil))
}

// ---
// ---
// ---

func setupDatabaseConnection() {
	mongoServers := os.Getenv("MONGO_SERVERS")
	
	// ---
	
	if mongoServers == "" {
		return
	}
	
	// ---
	
	Info("connecting to mongo databases at", mongoServers)
	
	// ---
	
	s, e := mgo.Dial(mongoServers)
	
	if e != nil {
		Fatal(e)
	}
	
	defer s.Close()
	
	// ---
	
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	
	if mongoDatabase == "" {
		Fatal("undefined mongo database")
	}
	
	// ---
	
	mongoCollection := os.Getenv("MONGO_COLLECTION")
	
	if mongoCollection == "" {
		Fatal("undefined mongo collection")
	}
	
	// ---
	
	G.C = s.DB(mongoDatabase).C(mongoCollection)
}

// ---
// ---
// ---

func init() {
	setupDatabaseConnection()
}

// ---
