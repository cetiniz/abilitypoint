package main

import (
	"encoding/json"
	"fmt"
	"github.com/Gurpartap/logrus-stack"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io/ioutil"
	"net/http"
	"os"
)

var base_url, user_name, password string

func init() {
	/* CONFIGURATION CODE */
	jsonFile, err := os.Open("../internal/config.json")
	if err != nil {
		logError(err, "Config File")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configuration map[string]interface{}
	json.Unmarshal([]byte(byteValue), &configuration)

	if str, ok := configuration["base_url"].(string); ok {
		base_url = str
	} else {
		fmt.Println("Can't start server because database url wasn't specified in config!")
		os.Exit(1)
	}
	if str, ok := configuration["base_url"].(string); ok {
		user_name = str
	} else {
		fmt.Println("Can't start server because database username wasn't specified in config!")
		os.Exit(1)
	}
	if str, ok := configuration["base_url"].(string); ok {
		password = str
	} else {
		fmt.Println("Can't start server because database password wasn't specified in config!")
		os.Exit(1)
	}

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

	if driver, err = neo4j.NewDriver("bolt://"+base_url+":7687", neo4j.BasicAuth(user_name, password, "")); err != nil {
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
