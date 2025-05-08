package tests

import (
	"github.com/saichler/gsql/go/gsql/interpreter"
	. "github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"testing"
)

func createQuery(query string) (*interpreter.Query, ifs.IResources, error) {
	r, _ := CreateResources(25000, 2, ifs.Trace_Level)
	r.Introspector().Inspect(&testtypes.TestProto{})
	q, e := interpreter.NewQuery(query, r)
	return q, r, e
}

func checkQuery(query string, expErr bool, t *testing.T) bool {
	q, _, e := createQuery(query)
	if e != nil && !expErr {
		Log.Fail(t, "Error creating query: ", e.Error())
		return false
	}
	if e == nil && expErr {
		Log.Fail(t, "Expected an error when creating a query")
		return false
	}
	if q == nil && e == nil {
		Log.Fail(t, "Query is nil")
		return false
	}
	return true
}

func checkMatch(query string, pb *testtypes.TestProto, expectMatch bool, t *testing.T) bool {
	q, _, e := createQuery(query)
	if e != nil {
		Log.Fail(t, e)
		return false
	}
	if !q.Match(pb) && expectMatch {
		Log.Fail(t, "Expected a match")
		return false
	}
	if q.Match(pb) && !expectMatch {
		Log.Fail(t, "Expected no match")
		return false
	}
	return true
}
