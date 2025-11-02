package storage

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

/*
Neo4jInstance handles interaction with a neo4j database
*/
type Neo4jInstance struct {
	Username      string
	Password      string
	URI           string
	DBName        string
	configOptions neo4j.ExecuteQueryConfigurationOption
	driver        neo4j.DriverWithContext
	ctx           context.Context
	debug         bool
}

/*
Init the instance facade
*/
func (neo *Neo4jInstance) Init() error {
	neo.ctx = context.Background()
	driver, err := neo4j.NewDriverWithContext(neo.URI, neo4j.BasicAuth(neo.Username, neo.Password, ""))
	if err != nil {
		return err
	}
	neo.driver = driver
	neo.configOptions = neo4j.ExecuteQueryWithDatabase(neo.DBName)
	return nil
}

/*
Verify the connection
*/
func (neo *Neo4jInstance) Verify() error {
	return neo.driver.VerifyConnectivity(neo.ctx)
}

/*
Clean all nodes and edges from the database
*/
func (neo *Neo4jInstance) Clean() error {
	neo.Execute(`match (a) -[r] -> () delete a, r`, map[string]any{})
	neo.Execute(`match (a) delete a`, map[string]any{})
	return nil
}

/*
Close the driver to the database
*/
func (neo *Neo4jInstance) Close() {
	_ = neo.driver.Close(neo.ctx)
}

/*
SetDebug on or off by boolean
*/
func (neo *Neo4jInstance) SetDebug(debug bool) {
	neo.debug = debug
}

/*
Execute the given query with the params
*/
func (neo *Neo4jInstance) Execute(query string, params map[string]any) {
	result, err := neo4j.ExecuteQuery(neo.ctx, neo.driver, query, params, neo4j.EagerResultTransformer, neo.configOptions)
	if err != nil {
		panic(err)
	}
	summary := result.Summary
	if neo.debug {
		fmt.Printf("Created %v nodes in %+v.\n", summary.Counters().NodesCreated(), summary.ResultAvailableAfter())
	}
}
