package flash

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

//Server model
type Server struct {
	*gin.Engine
	MetaData
	services   []interface{}
	middleware map[string]func(http.Handler) http.Handler
}

//Default will return default service engine
func Default() Server {
	return Server{
		Engine:     gin.Default(),
		services:   make([]interface{}, 0),
		middleware: make(map[string]func(http.Handler) http.Handler),
	}
}

//AddService will add new service in server
//service is a pointer to the struct of the api
func (s *Server) AddService(service interface{}) {
	if isStructAddress(service) {
		s.services = append(s.services, service)
	} else {
		panic("Expects an address of the Struct")
	}
}

//Start will start the server
func (s *Server) Start(port ...string) (err error) {
	s.setupServices()
	err = s.Run(port...)
	return
}

func (s *Server) setupServices() {
	for _, curService := range s.services {
		s.setupAPI(curService)
	}
}

//setAPI will set API metadata, middleware and handler
func (s *Server) setupAPI(service interface{}) {
	rootData := getMetaData(service, getStructName(Server{}))

	refObj := reflect.ValueOf(service).Elem()
	for i := 0; i < refObj.NumField(); i++ {
		field := refObj.Type().Field(i)
		if isAPI(field) {

			apiData := getMetaData(service, field.Name)
			metaData := getInOrderMetaData(apiData, rootData)

			if metaData.root == empty {
				metaData.root = strings.ToLower(refObj.Type().Name())
			}
			if metaData.url == empty {
				metaData.url = strings.ToLower(field.Name)
			}

			// method := getHTTPMethod(field)
			// handler := getHandler(service, firstCap(field.Name))
			s.setupHandler(service, field, metaData)
		}
	}
}

//setupHandler will setup Handler for the api
func (s *Server) setupHandler(service interface{}, field reflect.StructField, metaData MetaData) {
	switch field.Type.String() {
	case reflect.TypeOf(GET{}).String():
		s.GET(metaData.url, s.handler(service, field))
	case reflect.TypeOf(POST{}).String():
	case reflect.TypeOf(PUT{}).String():
	case reflect.TypeOf(PATCH{}).String():
	case reflect.TypeOf(DELETE{}).String():
	}
}

func (s *Server) handler(service interface{}, field reflect.StructField) func(c *gin.Context) {

	return func(c *gin.Context) {
		val := []reflect.Value{reflect.ValueOf(c)}
		reflect.ValueOf(service).MethodByName(firstCap(field.Name)).Call(val)
	}
}
