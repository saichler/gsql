package orm

import (
	"fmt"
	. "github.com/saichler/gsql/golang/parser"
	. "github.com/saichler/utils/golang"
	"strconv"
	"testing"
)

func TestQuery01(t *testing.T) {
	q, e := NewQuery("Select column1,column2 fRom table1 wHere 1=2")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	testTables(q, []string{"table1"}, t)
	testColumns(q, []string{"column1", "column2"}, t)
	testExpression(q, "(1=2)", t)
}

func TestQuery02(t *testing.T) {
	q, e := NewQuery("Select column1 fRom table1,table2 wHere 1=2")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	testTables(q, []string{"table1", "table2"}, t)
	testColumns(q, []string{"column1"}, t)
	testExpression(q, "(1=2)", t)
}

func TestQuery03(t *testing.T) {
	q, e := NewQuery("Select column1,column2 fRom table1,table2 wHere 1=2 AND 3=4")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	testTables(q, []string{"table1", "table2"}, t)
	testColumns(q, []string{"column1", "column2"}, t)
	testExpression(q, "(1=2 and 3=4)", t)
}

func TestQuery04(t *testing.T) {
	q, e := NewQuery("Select column1,column2 fRom table1,table2 wHere 1=2 AND 3  =  4")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	testTables(q, []string{"table1", "table2"}, t)
	testColumns(q, []string{"column1", "column2"}, t)
	testExpression(q, "(1=2 and 3=4)", t)
}

func TestQuery05(t *testing.T) {
	q, e := NewQuery("Select column1,column2 fRom table1,table2 wHere 1=2 AND 3  =  4 Or 5!=6")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	testTables(q, []string{"table1", "table2"}, t)
	testColumns(q, []string{"column1", "column2"}, t)
	testExpression(q, "(1=2 and 3=4 or 5!=6)", t)
}

func TestQuery06(t *testing.T) {
	q, e := NewQuery("Select column1,column2 fRom table1,table2 wHere 1=2 AND (3  =  4 Or 5!=6)")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	testTables(q, []string{"table1", "table2"}, t)
	testColumns(q, []string{"column1", "column2"}, t)
	testExpression(q, "(1=2) and ((3=4 or 5!=6))", t)
}

func TestQuery07(t *testing.T) {
	q, e := NewQuery("Select column1,column2 fRom table1,table2 wHere (1=2 or 3  =  4) And (5!=6 or 8<9) or 10<=12")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	testTables(q, []string{"table1", "table2"}, t)
	testColumns(q, []string{"column1", "column2"}, t)
	testExpression(q, "((1=2 or 3=4)) and ((5!=6 or 8<9)) or (10<=12)", t)
}

func TestQuery08(t *testing.T) {
	_, e := NewQuery("Select column1,column2 fRom table1,table2 wHere (1=2 or 3  =  4) And (5!=6 or 8<9 or 10<=12")
	if e == nil {
		Error("Expected failure.")
		t.Fail()
		return
	}
}

func TestQuery09(t *testing.T) {
	_, e := NewQuery("Select column1,column2 fRom table1,table2 wHere (1=2 or 3  =  4) And 5!=6 or 8<9) or 10<=12")
	if e == nil {
		Error("Expected failure.")
		t.Fail()
		return
	}
}

func TestQuery10(t *testing.T) {
	_, e := NewQuery("Select column1,column2 fRom table1,table2 wHere (1=2 or 3  =  4) Anf (5!=6 or 8<9) or 10<=12")
	if e == nil {
		Error("Expected failure.")
		t.Fail()
		return
	}
}

func TestQuery11(t *testing.T) {
	_, e := NewQuery("Select column1,column2 fRom table1,table2 wHere (1=2 or 3  =  4) And (5^6 or 8<9) or 10<=12")
	if e == nil {
		Error("Expected failure.")
		t.Fail()
		return
	}
}

func TestQuery12(t *testing.T) {
	q, e := NewQuery("Select column1,column2 fRom table1,table2 wHere (1=2 or 3  =  4) And (5!=6 or 8<9) or 10<=12 sort-by col1 page 7 limit 50 match-case descending")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	if !q.MatchCase() {
		t.Fail()
		Error("Expected Match Case to be true")
		return
	}
	if !q.Descending() {
		t.Fail()
		Error("Expected Descending to be true")
		return
	}
	if q.SortBy() != "col1" {
		t.Fail()
		Error("Expected sort-by to be col1")
		return
	}
	if q.Page() != 7 {
		t.Fail()
		Error("Expected page to be 7")
		return
	}
	if q.Limit() != 50 {
		t.Fail()
		Error("Expected kimit to be 50")
		return
	}

	testTables(q, []string{"table1", "table2"}, t)
	testColumns(q, []string{"column1", "column2"}, t)
	testExpression(q, "((1=2 or 3=4)) and ((5!=6 or 8<9)) or (10<=12)", t)
}

func TestVisualize(t *testing.T) {
	q, e := NewQuery("Select column1,column2 fRom table1,table2 wHere (1=2 or 3  =  4) And (5!=6 or 8<9) or 10<=12")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	fmt.Println(q.Where().String())
	fmt.Println(q.Where().Visualize(0))
}

func TestQuery(t *testing.T) {
	q, e := NewQuery("Select column1,column2 fRom table1 wHere 1=2 or ((3!=4 and 5<6) and 7>8) or ((9=10) and 11=12) ")
	if e != nil {
		Error(e)
		t.Fail()
		return
	}
	testTables(q, []string{"table1"}, t)
	testColumns(q, []string{"column1", "column2"}, t)
	testExpression(q, "(1=2) or (((3!=4 and 5<6)) and (7>8)) or (((9=10)) and (11=12))", t)
}

func testTables(q *Query, expected []string, t *testing.T) {
	if len(q.Tables()) != len(expected) {
		t.Fail()
		Error("Expected " + strconv.Itoa(len(expected)) + " tables but got " + strconv.Itoa(len(q.Tables())))
		return
	}
	for _, et := range expected {
		found := false
		for _, qt := range q.Tables() {
			if qt == et {
				found = true
				break
			}
		}
		if !found {
			t.Fail()
			Error("Expected table " + et + " but did not find it")
			return
		}
	}
}

func testColumns(q *Query, expected []string, t *testing.T) {
	if len(q.Columns()) != len(expected) {
		t.Fail()
		Error("Expected " + strconv.Itoa(len(expected)) + " columns but got " + strconv.Itoa(len(q.Columns())))
		return
	}
	for _, et := range expected {
		found := false
		for _, qc := range q.Columns() {
			if qc == et {
				found = true
				break
			}
		}
		if !found {
			t.Fail()
			Error("Expected column " + et + " but did not find it")
			return
		}
	}
}

func testExpression(q *Query, expected string, t *testing.T) {
	if q.Where().String() != expected {
		t.Fail()
		Error("Expected: " + expected)
		Error("But got : " + q.Where().String())
	}
}
