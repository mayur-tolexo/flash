package flash

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//TestService service struct
type TestService struct {
	Server `version:"1" root:"/test/"`
	ping   GET `url:"/ping"`
	ping2  GET `url:"/ping" version:"2"`
}

var router *Server

func init() {
	router = setupRoute()
}

//Ping api defination
func (*TestService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func setupRoute() (router *Server) {
	router = Default()
	router.AddService(&TestService{})
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
