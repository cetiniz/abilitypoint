package main

import (
	"fmt"
	"github.com/Gurpartap/logrus-stack"
	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"net/http"
	"os"
)

// Config contains all environment variables
type Config struct {
	BaseURL  string `env:"BASE_URL"`
	UserName string `env:"USER_NAME"`
	UserPass string `env:"USER_PASS"`
}

var cfg Config

func init() {
	// Configuration ENV Code
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading configuration from ENV")
	}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Failed to parse ENV")
		os.Exit(1)
	}

	fmt.Printf("Password: %s\n", cfg.UserPass)

	// Code for nicer looking logger
	log.AddHook(logrus_stack.StandardHook())
	log.SetOutput(os.Stderr)
	formatter := new(prefixed.TextFormatter)
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
	/* ******************  */

}

func main() {
	/* SETUP ROUTE HANDLER */
	r := mux.NewRouter()
	//r.HandleFunc("/relationship/{relationship_id}", relationshipHandler)
	http.Handle("/", r)
	r.HandleFunc("/api", YourHandler)
	r.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(":4000", r))
}

func fetchNeoGraph() []Edge {
	var (
		driver  neo4j.Driver
		session neo4j.Session
		result  neo4j.Result
		err     error
	)

	if driver, err = neo4j.NewDriver("bolt://"+cfg.BaseURL+":7687", neo4j.BasicAuth(cfg.UserName, cfg.UserPass, "")); err != nil {
		fmt.Println("ERROR")
	}
	// Used to destroy driver after calls
	defer driver.Close()

	if session, err = driver.Session(neo4j.AccessModeWrite); err != nil {
		fmt.Println("ERROR")
	}
	defer session.Close()

	result, err = session.Run("MATCH path = (n:Skill {name:'Integration by Parts'})-[:REQUIRES_UNDERSTANDING*0..2]->(j) WITH *, relationships(path) AS rels WITH [r IN rels | [STARTNODE(r), type(r), ENDNODE(r)]] AS steps, path UNWIND steps AS step RETURN DISTINCT step", nil)
	if err != nil {
		logError(err, "Neo4j")
	}

	var skills []Edge

	for result.Next() {
		currentNode := result.Record().GetByIndex(0).([]interface{})
		var From Skill
		var relationship string
		var To Skill
		err := mapstructure.Decode(currentNode[0].(neo4j.Node).Props(), &From)
		if err != nil {
			logError(err, "Decoding")
		}
		relationship, ok := currentNode[1].(string)
		if !ok {
			logError(err, "Decoding")
		}
		err = mapstructure.Decode(currentNode[2].(neo4j.Node).Props(), &To)
		if err != nil {
			logError(err, "Decoding")
		}
		var newEdge = Edge{
			From,
			relationship,
			To,
		}
		skills = append(skills, newEdge)
	}
	if err = result.Err(); err != nil {
		logError(err, "Neo4j")
	}
	return skills
}

func logError(err error, errType string) {
	log.WithFields(log.Fields{
		"Error_Message": err,
	}).Fatal(errType)
}

type Skill struct {
	Resources   []interface{}
	Images      []interface{}
	Name        string
	Description string
}

type Edge struct {
	From Skill
	Name string
	To   Skill
}

type MiddlewareFunc func(http.Handler) http.Handler

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
