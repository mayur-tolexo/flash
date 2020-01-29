package flash

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//Service1 service struct
type Service1 struct {
	Server          `version:"1" root:"/test/" prefix:"/"`
	ping            GET `url:"/ping"`
	ping2           GET `url:"/ping" version:"2"`
	extraParam      GET `url:"/invalid" v:"1"`
	ctxParamMissing GET `url:"/invalid" v:"2"`
}

//Service2 service struct
type Service2 struct {
	Server
	getPing    GET     `url:"/ping"`
	postPing   POST    `url:"/ping"`
	putPing    PUT     `url:"/ping"`
	patchPing  PATCH   `url:"/ping"`
	deletePing DELETE  `url:"/ping"`
	optionPing OPTIONS `url:"/ping"`
}

//Middlewares defined only on ping service endpoints
func (*Service2) Middlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{printMiddleware()}
}

func printMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		fmt.Println("Service Middleware Response time", time.Since(start).Seconds())
	}
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

func setupRoute() *Server {
	router := Default()
	router.AddService(&Service1{})
	router.AddService(&Service2{})
	return router
}

func TestServices(t *testing.T) {
	tc := []testCase{
		getPongRespTC("GET", "ping v1 success test", "/v1/test/ping"),
		get404TC("GET", "ping2 method not created so 404 check", "/v2/test/ping"),
		get404TC("GET", "extra param method 404 check", "/v1/test/invalid"),
		get404TC("GET", "ctx missing method 404 check", "/v2/test/invalid"),
		get404TC("GET", "ping2 method not created so 404 check", "/v2/test/ping"),
		getPongRespTC("GET", "service2 get ping check", "/service2/ping"),
		getPongRespTC("POST", "service2 post ping check", "/service2/ping"),
		getPongRespTC("PUT", "service2 put url check", "/service2/ping"),
		getPongRespTC("PATCH", "service2 patch url check", "/service2/ping"),
		getPongRespTC("DELETE", "service2 delete url check", "/service2/ping"),
		getPongRespTC("OPTIONS", "service2 option url check", "/service2/ping"),
	}
	router := setupRoute()
	runTestcase(t, router, tc)
}

func TestServer(t *testing.T) {
	router := New()
	router.AddService(&Service1{})
	port := ":7071"
	go router.Start(port)
	url := fmt.Sprintf("http://localhost%v/v1/test/ping", port)
	code, _, _ := getResp(url)
	assert := assert.New(t)
	assert.Equal(http.StatusOK, code)
}

func runTestcase(t *testing.T, router *Server, tc []testCase) {
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
	router := New()
	err := router.AddService(Service1{})
	assert.Error(t, err)
}

//Ping api defination
func (*Service1) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//ExtraParam api defination
func (*Service1) ExtraParam(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//CtxParamMissing api defination
func (*Service1) CtxParamMissing(data interface{}) {
}

//GetPing api defination
func (*Service2) GetPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//PostPing api defination
func (*Service2) PostPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//PutPing api defination
func (*Service2) PutPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//PatchPing api defination
func (*Service2) PatchPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//DeletePing api defination
func (*Service2) DeletePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//OptionPing api defination
func (*Service2) OptionPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func getResp(url string) (httpCode int, contentType string, content string) {
	client := &http.Client{}
	if req, err := http.NewRequest("GET", url, nil); err == nil {
		if resp, err := client.Do(req); err == nil {
			defer resp.Body.Close()
			if data, err := ioutil.ReadAll(resp.Body); err == nil {
				httpCode = resp.StatusCode
				contentType = resp.Header.Get("Content-Type")
				content = string(data)
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	return
}
