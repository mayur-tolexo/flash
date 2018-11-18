package flash

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

//isStructAddress will check the given input is pointer to struct
func isStructAddress(pt interface{}) (flag bool) {
	refObj := reflect.ValueOf(pt)
	if refObj.Kind() == reflect.Ptr && refObj.Elem().Kind() == reflect.Struct {
		flag = true
	}
	return
}

func getStructName(obj interface{}) string {
	return reflect.ValueOf(obj).Type().Name()
}

//coalesce will return first non-empty matching tag value for given variations of a key
func coalesce(tag reflect.StructTag, keys ...string) string {
	for _, key := range keys {
		if tag.Get(key) != "" {
			return tag.Get(key)
		}
	}
	return ""
}

//isAPI will check the service field an api or not
//field type should be api struct types
func isAPI(field reflect.StructField) (flag bool) {
	fType := field.Type.String()
	if field.Type.Kind() == reflect.Struct && strings.HasPrefix(fType, "flash") {
		flag = true
	}
	return
}

func firstCap(text string) string {
	out := []rune(text)
	if len(text) > 0 {
		out[0] = unicode.ToUpper(out[0])
	}
	return string(out)
}

func isSame(a, b reflect.Type) (flag bool) {
	if a.Kind() == b.Kind() && a.PkgPath() == b.PkgPath() && a.Name() == b.Name() {
		flag = true
	}
	return
}

func cleanURL(pieces ...string) string {

	var buffer bytes.Buffer

	// init the buffer to be a relative url
	buffer.WriteString("/")

	for _, p := range pieces {
		if p != "" && p != "-" {
			buffer.WriteString("/")
			buffer.WriteString(p)
		}
	}

	url := removeMultSlashes(buffer.String())
	//url = dropPrefix(url, "/")

	return url
}

var find *regexp.Regexp

func removeMultSlashes(inp string) string {
	if find == nil {
		find, _ = regexp.Compile("[\\/]+")
	}

	return find.ReplaceAllString(inp, "/")
}
