package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	graphql "github.com/playlyfe/go-graphql"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

var progress *graphql.Executor

func main() {

	var dbErr error
	db, dbErr = bolt.Open("./db/data.db", 0666,
		&bbolt.Options{NoSync: true})
	// above option necessary due to horror of:
	// https://github.com/etcd-io/bbolt/issues/149,
	// slow beyond usefuleness without.
	dbErr = ensureBuckets()
	if dbErr != nil {
		log.Fatal("cannot open database:", dbErr)
	}
	defer db.Close()
	log.Println("databse open")

	// commit data files
	// walk the data directories of json files & commit them to db
	dataPath := "./data/sif"
	err := filepath.Walk(dataPath, visitAndCommitSIF)
	if err != nil {
		log.Fatal("cannot load sif data", err)
	}
	dataPath = "./data/xapi"
	err = filepath.Walk(dataPath, visitAndCommitXAPI)
	if err != nil {
		log.Fatal("cannot load xapi data", err)
	}

	// construct gql resolvers & schema
	resolvers := buildResolvers()
	schema, err := buildSchema()
	if err != nil {
		log.Fatal("cannot load executor files: ", err)
	}

	var executorErr error
	progress, executorErr = graphql.NewExecutor(schema, "progress", "", resolvers)
	if executorErr != nil {
		log.Fatal("cannot create gql executor: ", executorErr)
	}

	// start the gql web server
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})) // allow cors requests during testing

	// entry point for javascript/css resources etc.
	e.Static("/", "public")

	// the graphql handlers
	e.POST("/graphql", gqlHandlerProgress)
	// e.POST("/search/graphql", gqlHandlerSearch)

	// run the server
	e.Logger.Fatal(e.Start(":1340"))

	// time.Sleep(time.Second * 5)

}

//
// wrapper type to capture graphql input
//
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

//
// the graphql handler routine for sylabus requests
//
func gqlHandlerProgress(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}

	result, err := progress.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	// log.Printf("result:\n\n%#v\n\n", result)

	return c.JSON(http.StatusOK, result)

}
