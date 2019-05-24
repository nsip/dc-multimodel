// schema.go

//
// handles loading of all of the gql schema files
// for use by the executor
//

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var schema []byte

func buildSchema() (string, error) {

	// walk the directory of json files & commit them to db
	root := "./gql"
	err := filepath.Walk(root, visitAndConcatenate)

	return string(schema), err

}

//
// called for each file the graphql directory
// builds a single string of all schema definitions
//
func visitAndConcatenate(path string, fi os.FileInfo, err error) error {

	//
	// make sure we only list .json files
	//
	if fi.Mode().IsRegular() && strings.HasSuffix(path, ".gql") {
		concatErr := concatGQL(path)
		if concatErr != nil {
			return concatErr
		}
		// log.Println("schema successfully added: ", path)
	}

	return nil
}

//
// reads the content of a grpahql idl file and appends it
// to the overall schema (string)
//
func concatGQL(path string) error {

	// read the file
	gqlFile, err := os.Open(path)
	if err != nil {
		return err
	}

	// retrieve bytes into json string
	gqlBytes, err := ioutil.ReadAll(gqlFile)
	if err != nil {
		return err
	}

	schema = append(schema, gqlBytes...)

	return nil
}
