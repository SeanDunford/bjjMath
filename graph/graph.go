package graph

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SeanDunford/bjjMath/scraper"
	"github.com/SeanDunford/simpleGraphGo/simplegraph"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

const (
	apple    = `{"name":"Apple Computer Company","type":["company","start-up"],"founded":"April 1, 1976","id":"1"}`
	woz      = `{"id":"2","name":"Steve Wozniak","type":["person","engineer","founder"]}`
	wozNick  = `{"name":"Steve Wozniak","type":["person","engineer","founder"],"nickname":"Woz","id":"2"}`
	jobs     = `{"id":"3","name":"Steve Jobs","type":["person","designer","founder"]}`
	wayne    = `{"id":"4","name":"Ronald Wayne","type":["person","administrator","founder"]}`
	markkula = `{"name":"Mike Markkula","type":["person","investor"]}`
	founded  = `{"action":"founded"}`
	invested = `{"action":"invested","equity":80000,"debt":170000}`
	divested = `{"action":"divested","amount":800,"date":"April 12, 1976"}`
)

func AddAthlete(athlete scraper.Athlete, dbFilePath string) {
	simplegraph.Initialize(dbFilePath)
	jsonAthlete, err := json.Marshal(athlete)
	if err != nil {
		fmt.Println(err)
		return
	}
	count, err := simplegraph.AddNode(athlete.Index, []byte(jsonAthlete), dbFilePath)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	fmt.Println(count)
}

func DoTheGraphTings(dbFilePath string) {
	simplegraph.Initialize(dbFilePath)
	_, _ = simplegraph.AddNode("1", []byte(apple), dbFilePath)
	_, _ = simplegraph.AddNode("2", []byte(woz), dbFilePath)
	_, _ = simplegraph.AddNode("3", []byte(jobs), dbFilePath)
	_, _ = simplegraph.AddNode("4", []byte(wayne), dbFilePath)
	_, _ = simplegraph.AddNode("5", []byte(markkula), dbFilePath)
	_, _ = simplegraph.AddNode("1", []byte(apple), dbFilePath)
	_, _ = simplegraph.AddNode("2", []byte(woz), dbFilePath)
	_, _ = simplegraph.ConnectNodesWithProperties("2", "1", []byte(founded), dbFilePath)
	_, _ = simplegraph.ConnectNodesWithProperties("3", "1", []byte(founded), dbFilePath)
	_, _ = simplegraph.ConnectNodesWithProperties("4", "1", []byte(founded), dbFilePath)
	_, _ = simplegraph.ConnectNodesWithProperties("5", "1", []byte(invested), dbFilePath)
	_, _ = simplegraph.ConnectNodesWithProperties("1", "4", []byte(divested), dbFilePath)
	_, _ = simplegraph.ConnectNodes("2", "3", dbFilePath)
	_, _ = simplegraph.FindNode("1", dbFilePath)
	_, _ = simplegraph.FindNode("7", dbFilePath)
	kvNameLike := simplegraph.GenerateWhereClause(&simplegraph.WhereClause{KeyValue: true, Key: "name", Predicate: "LIKE"})
	statement := simplegraph.GenerateSearchStatement(&simplegraph.SearchQuery{ResultColumn: "body", SearchClauses: []string{kvNameLike}})
	_, _ = simplegraph.FindNodes(statement, []string{"Steve%"}, dbFilePath)
	_ = simplegraph.UpdateNodeBody("2", wozNick, dbFilePath)
	_ = simplegraph.UpsertNode("1", apple, dbFilePath)
	_, _ = simplegraph.FindNode("2", dbFilePath)
	arrayType := simplegraph.GenerateWhereClause(&simplegraph.WhereClause{Tree: true, Predicate: "="})
	statement = simplegraph.GenerateSearchStatement(&simplegraph.SearchQuery{ResultColumn: "body", Tree: true, Key: "type", SearchClauses: []string{arrayType}})
	_, _ = simplegraph.FindNodes(statement, []string{"founder"}, dbFilePath)
	basicTraversal := simplegraph.GenerateTraversal(&simplegraph.Traversal{WithBodies: false, Inbound: true, Outbound: true})
	_, _ = simplegraph.TraverseFromTo("2", "3", basicTraversal, dbFilePath)
	basicTraversalInbound := simplegraph.GenerateTraversal(&simplegraph.Traversal{WithBodies: false, Inbound: true, Outbound: false})
	_, _ = simplegraph.TraverseFrom("5", basicTraversalInbound, dbFilePath)
	basicTraversalOutbound := simplegraph.GenerateTraversal(&simplegraph.Traversal{WithBodies: false, Inbound: false, Outbound: true})
	_, _ = simplegraph.TraverseFrom("5", basicTraversalOutbound, dbFilePath)
	_, _ = simplegraph.TraverseFrom("5", basicTraversal, dbFilePath)
	basicTraversalWithBodies := simplegraph.GenerateTraversal(&simplegraph.Traversal{WithBodies: true, Inbound: true, Outbound: true})
	_ = simplegraph.NodeData{Identifier: nil, Body: nil}
	_, _ = simplegraph.TraverseWithBodiesFromTo("2", "3", basicTraversalWithBodies, dbFilePath)
	_, _ = simplegraph.ConnectionsIn("1", dbFilePath)
	_ = []simplegraph.EdgeData{{Source: "1", Target: "4", Label: divested}}
	_, _ = simplegraph.ConnectionsOut("1", dbFilePath)
	_, _ = simplegraph.Connections("1", dbFilePath)
	_, _ = simplegraph.FindNode("2", dbFilePath)
	_, _ = simplegraph.FindNode("4", dbFilePath)

}
