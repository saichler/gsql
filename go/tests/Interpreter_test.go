package tests

import (
	. "github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/types/go/testtypes"
	"testing"
)

func TestQueryValidation(t *testing.T) {
	checkQuery("Select MyString fRom TeStproto wHere (MyString=hello world or (MyString=hello orm and myInt32=myvalue and mymodelslice=192*))",
		false, t)
}

func TestQueryMatch(t *testing.T) {
	checkQuery("Select MyString fRom testproto wHere (MyString=hello world or (MyString=hello orm and Myint32=myvalue and mymodelslice=192*))",
		false, t)
}

func TestMatchValue(t *testing.T) {
	q, _, e := createQuery("Select MyString fRom testproto wHere (MyString=hello world or (myString=hello orm and myint32=31 and mymodelslice.myString=192))")
	if e != nil {
		Log.Fail(t, e)
		return
	}

	node := CreateTestModelInstance(1)

	if q.Match(node) {
		Log.Fail(t, "Expected no match")
		return
	}

	node.MyString = "hello world"
	if !q.Match(node) {
		Log.Fail(t, "Expected a match")
		return
	}

	node.MyString = "hello orm"
	node.MyInt32 = 31
	node.MyModelSlice[0].MyString = "193"

	if q.Match(node) {
		Log.Fail(t, "Expected no match")
		return
	}

	node.MyString = "hello orm"
	node.MyInt32 = 31
	node.MyModelSlice[0].MyString = "192"
	if !q.Match(node) {
		Log.Fail(t, "Expected a match")
		return
	}

	if !checkMatch("Select myString fRom testproto wHere mymodelslice.mystring=192", node, true, t) {
		return
	}

	if !checkMatch("Select MyString fRom TestPRoto wHere mymodelslice.mYsTring=192 or mymodelslice.myString=193", node, true, t) {
		return
	}

	if !checkMatch("Select MyString fRom TestPRoto wHere mymodelslice.mYsTring=194 or mymodelslice.myString=193", node, false, t) {
		return
	}
}

func TestMultiMatchValue(t *testing.T) {
	node := CreateTestModelInstance(1)
	if !checkMatch("Select myString fRom testproto wHere mymodelslice.mystring=192 or mymodelslice.mystring=192", node, false, t) {
		return
	}
}

func TestDeepMatchMultiValueMap(t *testing.T) {
	node := CreateTestModelInstance(1)
	for _, v := range node.MyString2ModelMap {
		v.MyString = "Subnode6-0-index-0"
		break
	}
	if !checkMatch("Select myString fRom testproto wHere MyString2ModelMap.myString=Subnode6-0-index-0", node, true, t) {
		return
	}
}

func TestDeepMatchMultiValueMap2(t *testing.T) {
	node := CreateTestModelInstance(1)
	if !checkMatch("Select myString fRom testproto wHere MyString2ModelMap.myString=Subnode6-0-index-0", node, false, t) {
		return
	}
}

func TestDeepMatchMultiValueMap3(t *testing.T) {
	node := CreateTestModelInstance(1)
	node.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{}
	node.MyString2ModelMap["newone"].MySubs = make(map[string]*testtypes.TestProtoSubSub)
	node.MyString2ModelMap["newone"].MySubs["newone"] = &testtypes.TestProtoSubSub{MyString: "Subnode6-0-index-0"}
	if !checkMatch("Select myString fRom testproto wHere MyString2ModelMap.mysubs.myString=Subnode6-0-index-0", node, true, t) {
		return
	}
}

func TestDeepMatchMultiValueSlice(t *testing.T) {
	node := CreateTestModelInstance(1)
	pb := &testtypes.TestProtoSub{}
	pb.MySubs = make(map[string]*testtypes.TestProtoSubSub)
	pb.MySubs["newone"] = &testtypes.TestProtoSubSub{MyString: "Subnode6-0-index-0"}
	node.MyModelSlice = append(node.MyModelSlice, pb)
	if !checkMatch("Select myString fRom testproto wHere MyModelSlice.mysubs.myString=Subnode6-0-index-0", node, true, t) {
		return
	}
}
