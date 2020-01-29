package flash

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *Server

func init() {
	router = setupRoute()
}

//TestService service struct
type TestService struct {
	Server          `version:"1" root:"/test/" prefix:"/"`
	ping            GET `url:"/ping"`
	ping2           GET `url:"/ping" version:"2"`
	extraParam      GET `url:"/invalid" v:"1"`
	ctxParamMissing GET `url:"/invalid" v:"2"`
}

//TestService2 service struct
type TestService2 struct {
	Server
	getPing    GET     `url:"/ping"`
	postPing   POST    `url:"/ping"`
	putPing    PUT     `url:"/ping"`
	patchPing  PATCH   `url:"/ping"`
	deletePing DELETE  `url:"/ping"`
	optionPing OPTIONS `url:"/ping"`
}

//Middlewares defined only on ping service endpoints
func (*TestService2) Middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{serviceMiddleware()}
}

func serviceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		fmt.Println("Service Middleware Response time", time.Since(start).Seconds())
	}
}

func setupRoute() (router *Server) {
	router = Default()
	router.AddService(&TestService{})
	router.AddService(&TestService2{})
	return
}

type testCase struct {
	method   string
	name     string
	url      string
	expected string
	status   int
}

func getPongRespTC(method, name, url string) testCase {
	return testCase{
		method:   method,
		name:     name,
		url:      url,
		status:   http.StatusOK,
		expected: `{"message":"pong"}`,
	}
}

func get404TC(method, name, url string) testCase {
	return testCase{
		method:   method,
		name:     name,
		url:      url,
		status:   http.StatusNotFound,
		expected: `404 page not found`,
	}
}

func TestPingService(t *testing.T) {
	tc := []testCase{
		getPongRespTC("GET", "ping v1 success test", "/v1/test/ping"),
		get404TC("GET", "ping2 method not created so 404 check", "/v2/test/ping"),
		get404TC("GET", "extra param method 404 check", "/v1/test/invalid"),
		get404TC("GET", "ctx missing method 404 check", "/v2/test/invalid"),
		get404TC("GET", "ping2 method not created so 404 check", "/v2/test/ping"),
		getPongRespTC("GET", "service2 get ping check", "/testservice2/ping"),
		getPongRespTC("POST", "service2 post ping check", "/testservice2/ping"),
		getPongRespTC("PUT", "service2 put url check", "/testservice2/ping"),
		getPongRespTC("PATCH", "service2 patch url check", "/testservice2/ping"),
		getPongRespTC("DELETE", "service2 delete url check", "/testservice2/ping"),
		getPongRespTC("OPTIONS", "service2 option url check", "/testservice2/ping"),
	}

	for _, ctc := range tc {
		t.Run(ctc.name, func(t *testing.T) {
			req := createNilBodyReq(t, ctc.method, ctc.url)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert := assert.New(t)
			assert.Equal(ctc.status, w.Code)
			assert.Equal(ctc.expected, strings.TrimSuffix(w.Body.String(), "\n"))
		})
	}
}

func createNilBodyReq(t *testing.T, method, url string) (req *http.Request) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	return
}

func TestAddService(t *testing.T) {
	err := router.AddService(TestService{})
	assert.Error(t, err)
}

//Ping api defination
func (*TestService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//InvalidMethod api defination
func (*TestService) InvalidMethod(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//CtxParamMissing api defination
func (*TestService) CtxParamMissing(data interface{}) {
}

//GetPing api defination
func (*TestService2) GetPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//PostPing api defination
func (*TestService2) PostPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//PutPing api defination
func (*TestService2) PutPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//PatchPing api defination
func (*TestService2) PatchPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//DeletePing api defination
func (*TestService2) DeletePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//OptionPing api defination
func (*TestService2) OptionPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
