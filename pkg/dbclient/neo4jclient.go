package dbclient

import (
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"net/http"
)

var (
	session neo4j.Session
	result  neo4j.Result
	err     error
)

// INeo4jClient contains all attributes for the client
type INeo4jClient interface {
	OpenNeo4jClient(url string, user string, pass string)
	CloseNeo4jClient()
	RunQuery(queryString string) ([]interface{}, error)
}

// RunQuery executes a neo4j query with given string
func (nc *Neo4jClient) RunQuery(queryString string) ([]interface{}, error) {
	if session, err = nc.driver.Session(neo4j.AccessModeWrite); err != nil {
		fmt.Println("ERROR")
	}
	defer session.Close()

	var resultList []interface{}

	result, err = session.Run(queryString, nil)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		currentNode := result.Record().GetByIndex(0).([]interface{})
		resultList = append(resultList, currentNode)
	}
	if err = result.Err(); err != nil {
		return nil, err
	}
	return resultList, nil
}

// Neo4jClient driver encapsulation
type Neo4jClient struct {
	driver neo4j.Driver
}

// Struct to hold statements
type neo4jStatements struct {
	statements []neo4jStatement
}

// Struct to hold single statement
type neo4jStatement struct {
	statement string
}

type neo4jResponse struct {
	results []neo4jResult
	errors  []neo4jError
}

type neo4jResult struct {
	columns []string
	data    []neo4jRow
}

type neo4jError struct {
}

type neo4jRow struct {
	row  []string
	meta []string
}

func (nc *Neo4jClient) executeSingleTransaction(query string) {
	var payload neo4jStatements
	payload.statements = make([]neo4jStatement, 0)
	payload.statements = append(payload.statements, neo4jStatement{statement: query})

	marshalledPayload, err := json.Marshal(payload)

	res, err := http.Post(nc.url+":7474:db/data/transaction/commit", "application/json", marshalledPayload)

	var neo4jRes neo4jResponse
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&neo4jRes)
}

// OpenNeo4jClient spawns client instance
func (nc *Neo4jClient) OpenNeo4jClient(url string, user string, pass string) {
	if nc.driver, err = neo4j.NewDriver("bolt://"+url+":7687", neo4j.BasicAuth(user, pass, "")); err != nil {
		fmt.Println("ERROR")
	}
}

// CloseNeo4jClient closes Bolt Driver
func (nc *Neo4jClient) CloseNeo4jClient() {
	defer nc.driver.Close()
}
