package jsonnet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUuidV5(t *testing.T) {
	var tests = []struct {
		input          []interface{}
		expectedOutput interface{}
		expectedErr    string
	}{
		{nil, "", "namespace and name must be provided"},
		{[]interface{}{""}, "", "namespace and name must be provided"},
		{[]interface{}{1, ""}, "", "namespace must be a string"},
		{[]interface{}{"", ""}, "", "namespace must be a UUID"},
		{[]interface{}{"dns", 1}, "", "name must be a string"},
		{[]interface{}{"dns", "www.terraform.io"}, "a5008fae-b28c-5ba5-96cd-82b4c53552d6", ""},
		{[]interface{}{"url", "https://www.terraform.io/"}, "9db6f67c-dd95-5ea0-aa5b-e70e5c5f7cf5", ""},
		{[]interface{}{"oid", "1.3.6.1.4"}, "af9d40a5-7a36-5c07-b23a-851cd99fbfa5", ""},
		{[]interface{}{"x500", "CN=Example,C=GB"}, "84e09961-4aa4-57f8-95b7-03edb1073253", ""},
		{[]interface{}{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", "www.terraform.io"}, "a5008fae-b28c-5ba5-96cd-82b4c53552d6", ""},
		{[]interface{}{"743ac3c0-3bf7-4a5b-9e6c-59360447c757", "LIBS:diskfont.library"}, "ede1a974-df7e-5f17-84b9-76208818b2c8", ""},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("uuidv5-%d", i), func(t *testing.T) {
			ret, err := UuidV5().Func(test.input)
			if test.expectedErr == "" {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutput, ret)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, test.expectedErr, err.Error())
			}
		})
	}
}
