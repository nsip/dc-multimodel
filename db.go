// db.go

//
// captures all db interactions
//

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tidwall/gjson"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

//
// make sure key-store buckets are in place
// to support our transactiosn and indexes.
//
//
func ensureBuckets() error {

	buckets := []string{
		"TeachingGroup",
		"StaffPersonal",
		"GradingAssignment",
		"StudentPersonal",
		"StudentAttendanceTimeList",
		"xApi",
		"TeachersByName"}

	var dbErr error
	for _, name := range buckets {
		dbErr = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return fmt.Errorf("create bucket failed: %s", err)
			}
			return nil
		})
	}

	return dbErr

}

func visitAndCommitSIF(path string, fi os.FileInfo, err error) error {

	//
	// make sure we only list .json files
	//
	if fi.Mode().IsRegular() && strings.HasSuffix(path, ".json") {
		commitErr := commitSIF(path)
		if commitErr != nil {
			return commitErr
		}
		// log.Println("successfully committed: ", path)
	}

	return nil
}

func commitSIF(path string) error {

	// read the file
	sifFile, err := os.Open(path)
	if err != nil {
		return err
	}

	// retrieve bytes into json string
	sifBytes, err := ioutil.ReadAll(sifFile)
	if err != nil {
		return err
	}
	sifContent := gjson.ParseBytes(sifBytes)
	if !gjson.Valid(sifContent.String()) {
		return fmt.Errorf("could not read valid json from file: %s", path)
	}

	var txErr error
	sifContent.ForEach(func(key, value gjson.Result) bool {
		// only one named item in the map for SIF objects:
		for k, v := range value.Map() {
			// log.Println("type: ", k)
			// log.Printf("\ntype:%s\nentry:\n%s\n\n", k, v.Raw)
			switch k {
			case "TeachingGroup":
				tgroupid := v.Get("RefId")
				teacherid := v.Get("TeacherList.TeachingGroupTeacher.0.StaffPersonalRefId")
				dbKey := fmt.Sprintf("%s:%s", teacherid, tgroupid)
				txErr = commit(k, dbKey, v.Raw)
				txErr = commit(k, tgroupid.String(), v.Raw)
			case "StaffPersonal":
				teacherid := v.Get("RefId")
				txErr = commit(k, teacherid.String(), v.Raw)
				// convenience lookup if no id is known
				firstName := v.Get("PersonInfo.Name.GivenName")
				lastName := v.Get("PersonInfo.Name.FamilyName")
				dbKey := fmt.Sprintf("%s %s", firstName, lastName)
				txErr = commit("TeachersByName", dbKey, teacherid.String())
			case "GradingAssignment":
				tgroupid := v.Get("TeachingGroupRefId")
				task := v.Get("DetailedDescriptionURL")
				dbKey := fmt.Sprintf("%s:%s", tgroupid, task)
				txErr = commit(k, dbKey, v.Raw)
			case "StudentAttendanceTimeList":
				spid := v.Get("StudentPersonalRefId")
				tlid := v.Get("RefId")
				dbKey := fmt.Sprintf("%s:%s", spid, tlid)
				txErr = commit(k, dbKey, v.Raw)
			}
		}
		return true // keep iterating
	})

	return txErr
}

func visitAndCommitXAPI(path string, fi os.FileInfo, err error) error {

	//
	// make sure we only list .json files
	//
	if fi.Mode().IsRegular() && strings.HasSuffix(path, ".json") {
		commitErr := commitXAPI(path)
		if commitErr != nil {
			return commitErr
		}
		// log.Println("successfully committed: ", path)
	}

	return nil
}

func commitXAPI(path string) error {

	// read the file
	xapiFile, err := os.Open(path)
	if err != nil {
		return err
	}

	// retrieve bytes into json string
	xapiBytes, err := ioutil.ReadAll(xapiFile)
	if err != nil {
		return err
	}
	xapiContent := gjson.ParseBytes(xapiBytes)
	if !gjson.Valid(xapiContent.String()) {
		return fmt.Errorf("could not read valid json from file: %s", path)
	}

	var txErr error
	xapiContent.ForEach(func(key, value gjson.Result) bool {
		// log.Printf("xapi-statement:\n%v\n\n", value.Map())
		name := value.Get("actor.name")
		task := value.Get("object.id")
		id := value.Get("id")
		dbKeu := fmt.Sprintf("%s:%s:%s", name, task, id)
		txErr = commit("xApi", dbKeu, value.Raw)
		return true // keep iterating
	})

	return txErr
}

func commit(bucket, key, value string) error {

	// log.Printf("\ntype: %s\nkey:%s\njson:\n%s\n\n", bucket, key, value)

	dbErr := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})

	return dbErr
}
