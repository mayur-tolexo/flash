package flash

import (
	"fmt"
	"net/http"
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

//Ping api defination
func (*TestService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func TestServer(t *testing.T) {
	router := Default()
	router.AddService(&TestService{})
	port := ":7071"
	go router.Start(port)
	url := fmt.Sprintf("http://localhost%v/v1/test/ping", port)
	code, _, data := getResp(url)
	assert := assert.New(t)
	assert.Equal(http.StatusOK, code)
	assert.Equal(`{"message":"pong"}`, data)

	url = fmt.Sprintf("http://localhost%v/v2/test/ping", port)
	code, _, _ = getResp(url)
	assert.Equal(http.StatusNotFound, code)
}
