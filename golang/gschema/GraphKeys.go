package gschema

import (
	"bytes"
	"strings"
)

type GraphKeys struct {
	keys map[string]string
}

func (graphKeys *GraphKeys) Key(path string) string {
	return graphKeys.keys[path]
}

func (graphKeys *GraphKeys) Strings() string {
	buff := bytes.Buffer{}
	for k, v := range graphKeys.keys {
		buff.WriteString("K=")
		buff.WriteString(k)
		buff.WriteString(" V=")
		buff.WriteString(v)
		buff.WriteString("\n")
	}
	return buff.String()
}

func NewGraphKeys(id string) (*GraphKeys, string) {
	id = RemoveUnderScore(id)
	path := TrimAndLowerNoKeys(id)
	gkeys := &GraphKeys{}
	gkeys.keys = make(map[string]string)
	from := strings.Index(path, "[")
	for from != -1 {
		to := strings.Index(path, "]")
		prefix := path[0:from]
		suffix := path[to+1:]
		key := path[from+1 : to]
		gkeys.keys[prefix] = key
		buff := bytes.Buffer{}
		buff.WriteString(prefix)
		buff.WriteString(suffix)
		path = buff.String()
		from = strings.Index(path, "[")
	}
	return gkeys, path
}

func RemoveUnderScore(id string) string {
	buff := bytes.Buffer{}
	keyOpen := false
	for _, c := range id {
		char := string(c)
		if char == "[" {
			keyOpen = true
		} else if char == "]" {
			keyOpen = false
		}
		if char != "_" || keyOpen {
			buff.WriteString(char)
		}
	}
	return buff.String()
}

func TrimAndLowerNoKeys(sql string) string {
	buff := bytes.Buffer{}
	sql = strings.TrimSpace(sql)
	keyOpen := false
	for _, c := range sql {
		if c == '[' {
			keyOpen = true
		} else if c == ']' {
			keyOpen = false
		}
		if !keyOpen {
			buff.WriteString(strings.ToLower(string(c)))
		} else {
			buff.WriteString(string(c))
		}
	}
	return buff.String()
}
