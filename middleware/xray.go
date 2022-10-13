package middleware

import (
	"net/http"
	"strings"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gin-gonic/gin"
)

// X-Rayでトレースを行うためのミドルウェア
func XrayMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, seg := xray.BeginSegment(c, c.Request.Method+" "+c.Request.URL.Path)
		c.Request = c.Request.WithContext(ctx)

		captureRequestData(c, seg)
		c.Next()
		captureResponseData(c, seg)

		seg.Close(nil)
	}
}

// リクエストデータのキャプチャ
func captureRequestData(c *gin.Context, seg *xray.Segment) {
	seg.Lock()
	defer seg.Unlock()

	r := c.Request

	requestData := seg.GetHTTP().GetRequest()
	requestData.Method = r.Method
	requestData.URL = r.URL.String()
	requestData.XForwardedFor = hasXForwardedFor(r)
	requestData.ClientIP = getClientIP(r)
	requestData.UserAgent = r.UserAgent()
}

// レスポンスデータのキャプチャ
func captureResponseData(c *gin.Context, seg *xray.Segment) {
	seg.Lock()
	defer seg.Unlock()

	statusCode := c.Writer.Status()

	responseData := seg.GetHTTP().GetResponse()
	responseData.Status = statusCode
	responseData.ContentLength = c.Writer.Size()

	if statusCode >= 400 && statusCode < 500 {
		seg.Error = true
	}
	if statusCode == 429 {
		seg.Throttle = true
	}
	if statusCode >= 500 && statusCode < 600 {
		seg.Fault = true
	}
}

// X-Forwarded-Forの有無チェック
func hasXForwardedFor(r *http.Request) bool {
	return r.Header.Get("X-Forwarded-For") != ""
}

// クライアントのIP取得
func getClientIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		return strings.TrimSpace(strings.Split(forwardedFor, ",")[0])
	}
	return r.RemoteAddr
}
