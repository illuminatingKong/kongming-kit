package webapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/illuminatingKong/kongming-kit/http/corehandler"
	"github.com/illuminatingKong/kongming-kit/http/webapi"
	"net/http"
	"testing"
)

var respSuccess = webapi.RespDefault

func TestWebHttpApiResponseHandler(t *testing.T) {
	router := gin.New()
	router.GET("/hello", func(c *gin.Context) {
		ctx := corehandler.NewContext(c)
		c.Set("requestID", "123")
		defer func() { corehandler.WebHttpApiResponse(c, ctx) }()
		ctx.Resp = respSuccess().SetData(map[string]interface{}{"hello": "kongming"}).
			SetMessage("say hello").SetTotal(20)
	})
	req, _ := http.NewRequest("GET", "/hello", nil)
	resp, _ := http.DefaultClient.Do(req)
	fmt.Println(resp)
}
