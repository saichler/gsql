package introspector

import (
	"github.com/saichler/gsql/golang/gschema"
	. "github.com/saichler/utils/golang"
	"reflect"
	"strconv"
	"strings"
)

type Annotation string

const (
	Title        Annotation = "Title"
	Size         Annotation = "Size"
	Mask         Annotation = "Mask"
	Ignore       Annotation = "Ignore"
	PrimaryKey   Annotation = "PrimaryKey"
	UniqueKey    Annotation = "UniqueKey"
	NonUniqueKey Annotation = "NonUniqueKey"
)

type Column struct {
	table    *Table
	field    reflect.StructField
	metaData *ColumnMetaData
}

func (c *Column) MetaData() *ColumnMetaData {
	return c.metaData
}

func (c *Column) inspect(parent *gschema.GraphSchemaNode) {
	c.parseMetaData()
	if isStruct(c.field.Type) {
		strct := getStruct(c.field.Type)
		c.metaData.columnTableName = strct.Name()
		schemaNode := c.table.introspector.graphSchema.RegisterStruct(c.Name(), parent, strct)
		if schemaNode != nil {
			c.table.introspector.introspect(strct, schemaNode)
		}
	}
}

func (c *Column) Type() reflect.Type {
	return c.field.Type
}

func (c *Column) Name() string {
	return c.field.Name
}

func (c *Column) Table() *Table {
	return c.table
}

func isStruct(typ reflect.Type) bool {
	if typ.Kind() == reflect.Struct {
		return true
	} else if typ.Kind() == reflect.Ptr {
		return isStruct(typ.Elem())
	} else if typ.Kind() == reflect.Slice {
		return isStruct(typ.Elem())
	} else if typ.Kind() == reflect.Map {
		return isStruct(typ.Elem())
	}
	return false
}

func getStruct(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Struct {
		return typ
	} else if typ.Kind() == reflect.Ptr {
		return getStruct(typ.Elem())
	} else if typ.Kind() == reflect.Slice {
		return getStruct(typ.Elem())
	} else if typ.Kind() == reflect.Map {
		return getStruct(typ.Elem())
	}
	return nil
}

func (c *Column) parseMetaData() {
	if c.metaData == nil {
		c.metaData = &ColumnMetaData{}
	}
	c.metaData.title = c.field.Name
	c.metaData.size = 128
	tags := string(c.field.Tag)
	annTags := c.table.introspector.AsTags(c.table.Name() + "." + c.Name())
	tags+=" "+ annTags
	if tags == "" {
		return
	}
	splits := strings.Split(tags, " ")
	for _, tag := range splits {
		c.getTag(tag)
	}
}

func (c *Column) getTag(tag string) {
	if strings.TrimSpace(tag) == "" {
		return
	}
	index := strings.Index(tag, "=")
	if index == -1 {
		return
	}
	name := Annotation(tag[0:index])
	value := tag[index+1:]
	if name == Title {
		c.metaData.title = value
	} else if name == Size {
		val, err := strconv.Atoi(value)
		if err != nil {
			Error("Unable to parse field size from:" + value + " in field:" + c.field.Name)
		} else {
			c.metaData.size = val
		}
	} else if name == Mask {
		c.metaData.mask = true
	} else if name == Ignore {
		c.metaData.ignore = true
	} else if name == PrimaryKey {
		c.metaData.primaryKey = value
	} else if name == UniqueKey {
		c.metaData.uniqueKeys = value
	} else if name == NonUniqueKey {
		c.metaData.nonUniqueKeys = value
	}
}
