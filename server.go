package flash

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

//Server model
type Server struct {
	*gin.Engine
	MetaData
	services []interface{}
}

//Default will return default service engine
func Default() *Server {
	return &Server{
		Engine:   gin.Default(),
		services: make([]interface{}, 0),
	}
}

//New will return new service engine
func New() *Server {
	return &Server{
		Engine:   gin.New(),
		services: make([]interface{}, 0),
	}
}

//AddService will add new service in server
//service is a pointer to the struct of the api
func (s *Server) AddService(service interface{}) (err error) {
	if isStructAddress(service) {
		s.services = append(s.services, service)
		s.setupAPI(service)
	} else {
		err = fmt.Errorf("Expects an address of the Struct")
	}
	return
}

//Start will start the server
func (s *Server) Start(port ...string) (err error) {
	err = s.Run(port...)
	return
}

//setAPI will set API metadata, middleware and handler
func (s *Server) setupAPI(service interface{}) {
	rootData := getMetaData(service, getStructName(Server{}))

	refObj := reflect.ValueOf(service).Elem()
	grp := s.Group("/")
	grp.Use(getServiceMiddlewares(service)...)

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

			method := getMethod(service, field)
			s.setupHandler(grp, method, metaData)
		}
	}
}

//setupHandler will setup Handler for the api
func (s *Server) setupHandler(grp *gin.RouterGroup, method Method, metaData MetaData) {
	if handler, exists := method.getHandler(); exists {
		var version string
		if metaData.version != "" {
			version = "v" + metaData.version
		}
		url := cleanURL(metaData.prefix, version, metaData.root, metaData.url)
		switch method.methodType {
		case reflect.TypeOf(GET{}).String():
			grp.GET(url, handler)
		case reflect.TypeOf(POST{}).String():
			grp.POST(url, handler)
		case reflect.TypeOf(PUT{}).String():
			grp.PUT(url, handler)
		case reflect.TypeOf(PATCH{}).String():
			grp.PATCH(url, handler)
		case reflect.TypeOf(DELETE{}).String():
			grp.DELETE(url, handler)
		case reflect.TypeOf(OPTIONS{}).String():
			grp.OPTIONS(url, handler)
		}
	}
}
