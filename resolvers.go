package main

import (
	"bytes"
	"fmt"
	"strconv"

	graphql "github.com/playlyfe/go-graphql"
	"github.com/tidwall/gjson"
	"go.etcd.io/bbolt"
)

func buildResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}

	//
	// returns the list of known teachers
	//
	resolvers["progress/teachers"] = func(params *graphql.ResolveParams) (interface{}, error) {

		teachers := make([]interface{}, 0)

		dbErr := db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte("StaffPersonal"))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				// fmt.Printf("key=%s, value=%s\n", k, v)
				sp := gjson.ParseBytes(v)
				teachers = append(teachers, sp.Value())
			}
			return nil
		})
		return teachers, dbErr
	}

	//
	// for a given teacher id returns their teaching groups
	//
	resolvers["progress/my_teaching_groups"] = func(params *graphql.ResolveParams) (interface{}, error) {

		tgroups := make([]interface{}, 0)

		id := fmt.Sprintf("%s", params.Args["teacher_id"])

		dbErr := db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte("TeachingGroup"))
			c := b.Cursor()
			for k, v := c.Seek([]byte(id)); k != nil && bytes.HasPrefix(k, []byte(id)); k, v = c.Next() {
				// fmt.Printf("key=%s, value=%s\n", k, v)
				sp := gjson.ParseBytes(v)
				tgroups = append(tgroups, sp.Value())
			}
			return nil
		})
		return tgroups, dbErr
	}

	//
	// for a given teaching group id, returns the xapi assessment event records for
	// each student in the teaching group.
	//
	resolvers["progress/reports"] = func(params *graphql.ResolveParams) (interface{}, error) {

		reports := []interface{}{}
		assignment_progress_report := map[string]interface{}{}

		// retrieve tg-id from query
		tgid := fmt.Sprintf("%s", params.Args["teaching_group_id"])

		// retrieve linked data from db
		dbErr := db.View(func(tx *bbolt.Tx) error {

			// get the teaching group
			b := tx.Bucket([]byte("TeachingGroup"))
			tg := gjson.ParseBytes(b.Get([]byte(tgid)))
			assignment_progress_report["teaching_group_name"] = tg.Get("ShortName")

			// get the related grading assignments
			assignemnts := []gjson.Result{}
			b = tx.Bucket([]byte("GradingAssignment"))
			c := b.Cursor()
			for k, v := c.Seek([]byte(tgid)); k != nil && bytes.HasPrefix(k, []byte(tgid)); k, v = c.Next() {
				ga := gjson.ParseBytes(v)
				assignemnts = append(assignemnts, ga)
			}

			// iterate assignements, and find xapi events relating to
			// the student and grading assignement
			results := []interface{}{}
			for _, ga := range assignemnts {

				taskName := fmt.Sprintf("%s", ga.Get("DetailedDescriptionURL"))
				assignemnt_results := map[string]interface{}{}

				student_results := []map[string]interface{}{}
				for _, student := range tg.Get("StudentList.TeachingGroupStudent").Array() {

					fullname := fmt.Sprintf("%s %s", student.Get("Name.GivenName"), student.Get("Name.FamilyName"))
					prefix := fmt.Sprintf("%s:%s", fullname, taskName)
					b = tx.Bucket([]byte("xApi"))
					c = b.Cursor()
					for k, v := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, v = c.Next() {

						xapievent := gjson.ParseBytes(v)
						// get the shorthand name for the assignment - easier to display & read
						assignemnt_results["assignment_name"] = xapievent.Get("object.definition.name").String()
						result := map[string]interface{}{"student_name": fullname, "result_event": xapievent.Value()}
						student_results = append(student_results, result)
						// get the attendnace summary
						bkt := tx.Bucket([]byte("StudentAttendanceTimeList"))
						crsr := bkt.Cursor()
						att_prefix := fmt.Sprintf("%s:", student.Get("StudentPersonalRefId"))
						missed_days := 0
						for k, v := crsr.Seek([]byte(att_prefix)); k != nil && bytes.HasPrefix(k, []byte(att_prefix)); k, v = crsr.Next() {
							_ = v
							missed_days++
						}
						result["absence_days"] = strconv.Itoa(missed_days)
					}
				}
				assignemnt_results["student_results"] = student_results
				results = append(results, assignemnt_results)
			}
			assignment_progress_report["assignment_results"] = results
			reports = append(reports, assignment_progress_report)
			return nil
		})
		return reports, dbErr

	}

	return resolvers
}
