package jsonnet

import (
	"fmt"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/google/uuid"
)

var uuidV5Namespaces = map[string]string{
	"dns":  "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	"url":  "6ba7b811-9dad-11d1-80b4-00c04fd430c8",
	"oid":  "6ba7b812-9dad-11d1-80b4-00c04fd430c8",
	"x500": "6ba7b814-9dad-11d1-80b4-00c04fd430c8",
}

func UuidV5() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "uuidv5",
		Params: ast.Identifiers{"namespace", "name"},
		Func: func(input []any) (any, error) {
			if len(input) != 2 {
				return nil, fmt.Errorf("namespace and name must be provided")
			}
			namespace, ok := input[0].(string)
			if !ok {
				return nil, fmt.Errorf("namespace must be a string")
			}
			namespaceString, ok := uuidV5Namespaces[namespace]
			if !ok {
				namespaceString = namespace
			}
			namespaceUuid, err := uuid.Parse(namespaceString)
			if err != nil {
				return nil, fmt.Errorf("namespace must be a UUID")
			}
			name, ok := input[1].(string)
			if !ok {
				return nil, fmt.Errorf("name must be a string")
			}
			res := uuid.NewSHA1(namespaceUuid, []byte(name))
			return res.String(), nil
		},
	}
}
