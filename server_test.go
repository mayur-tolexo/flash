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

func TestPingService(t *testing.T) {
	tc := []struct {
		name     string
		url      string
		expected string
		status   int
	}{
		{
			name:     "ping v1 success test",
			url:      "/v1/test/ping",
			status:   http.StatusOK,
			expected: `{"message":"pong"}`,
		},
		{
			name:     "ping2 method not created so 404 check",
			url:      "/v2/test/ping",
			status:   http.StatusNotFound,
			expected: `404 page not found`,
		},
	}

	for _, ctc := range tc {
		t.Run(ctc.name, func(t *testing.T) {
			req := createGetReq(t, ctc.url)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert := assert.New(t)
			assert.Equal(ctc.status, w.Code)
			assert.Equal(ctc.expected, strings.TrimSuffix(w.Body.String(), "\n"))
		})
	}
}

func createGetReq(t *testing.T, url string) (req *http.Request) {
	req, err := http.NewRequest("GET", url, nil)
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
