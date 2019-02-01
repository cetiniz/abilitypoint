package service

import (
	"encoding/json"
	"fmt"
	"github.com/cetiniz/abilitypoint/pkg/dbclient"
	"github.com/cetiniz/abilitypoint/pkg/model"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"net/http"
)

// DBClient is an instance of the neo4j db
var DBClient dbclient.INeo4jClient

func fetchGraph(w http.ResponseWriter, r *http.Request) {
	result, err := DBClient.RunQuery("MATCH path = (n:Skill {name:'Integration by Parts'})-[:REQUIRES_UNDERSTANDING*0..2]->(j) WITH *, relationships(path) AS rels WITH [r IN rels | [STARTNODE(r), type(r), ENDNODE(r)]] AS steps, path UNWIND steps AS step RETURN DISTINCT step")
	if err != nil {
		fmt.Println(err)
	}

	var skills []model.Edge

	for _, node := range result {
		var (
			From         model.Skill
			relationship string
			To           model.Skill
		)
		err := mapstructure.Decode(node.([]interface{})[0].(neo4j.Node).Props(), &From)
		if err != nil {
			fmt.Println(err, "Decoding")
		}
		relationship, ok := node.([]interface{})[1].(string)
		if !ok {
			fmt.Println(err, "Decoding")
		}
		err = mapstructure.Decode(node.([]interface{})[2].(neo4j.Node).Props(), &To)
		if err != nil {
			fmt.Println(err, "Decoding")
		}
		var newEdge = model.Edge{
			From,
			relationship,
			To,
		}
		skills = append(skills, newEdge)
	}

	js, err := json.Marshal(skills)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
