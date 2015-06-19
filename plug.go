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

func Getenv(name string) (string) {
	return os.Getenv(name)
}

// ---

func GetenvF(name string) (string) {
	value := os.Getenv(name)
	
	// ---
	
	if value == "" {
		Fatal("undefined ", name)
	}
	
	return value
}

func GetenvD(name string, defaultValue string) (string) {
	value := os.Getenv(name)
	
	// ---
	
	if value == "" {
		value = defaultValue
	}
	
	return defaultValue
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
	
	address := GetenvD("ADDRESS", ":8080")
	
	// ---
	
	Info("starting server at", address)
	
	// ---
	
	Fatal(http.ListenAndServe(address, nil))
}

// ---
// ---
// ---

func setupDatabaseConnection() {
	mongoServers := Getenv("MONGO_SERVERS")
	
	// ---
	
	if mongoServers == "" {
		return
	}
	
	// ---
	
	mongoDatabase := GetenvF("MONGO_DATABASE")
	mongoCollection := GetenvF("MONGO_COLLECTION")
	
	// ---
	
	Info("connecting to mongo databases at", mongoServers)
	
	// ---
	
	s, e := mgo.Dial(mongoServers)
	
	if e != nil {
		Fatal(e)
	}
	
	defer s.Close()
	
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
