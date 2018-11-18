package flash

import (
	"reflect"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

//Method need to call by the api
type Method struct {
	address    interface{}
	methodType string
	name       string
	params     []reflect.Value
	exists     bool
}

func getMethod(service interface{}, field reflect.StructField) Method {
	name := firstCap(field.Name)
	serviceName := reflect.TypeOf(service).Elem().Name()
	m, exists := reflect.TypeOf(service).MethodByName(name)
	if exists {
		d := color.New(color.FgHiBlue)
		if m.Type.NumIn() != 2 {
			d.Printf("[WARNING]: Need only one input of *gin.Context type in %v method of %v service\n\n", name, serviceName)
			exists = false
		} else {
			in := m.Type.In(1)
			if in.Kind() != reflect.Ptr || !isSame(in.Elem(), reflect.TypeOf(gin.Context{})) {
				d.Printf("*gin.Context missing in %v method of %v service\n\n", name, serviceName)
				exists = false
			}
		}
	}
	return Method{
		address:    service,
		methodType: field.Type.String(),
		name:       name,
		params:     make([]reflect.Value, 1),
		exists:     exists,
	}
}

func (m Method) getHandler() (func(c *gin.Context), bool) {
	return func(c *gin.Context) {
		m.params[0] = reflect.ValueOf(c)
		reflect.ValueOf(m.address).MethodByName(m.name).Call(m.params)
	}, m.exists
}
