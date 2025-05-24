package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/WangWilly/swap-estimation/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/smartystreets/goconvey/convey"
)

////////////////////////////////////////////////////////////////////////////////

// Controller defines the interface that controllers must implement
// to be used with the test HTTP server
type Controller interface {
	RegisterRoutes(router *gin.Engine)
}

////////////////////////////////////////////////////////////////////////////////

type TestHttpServer struct {
	Server *httptest.Server
	client *http.Client
}

func NewTestHttpServer(controller Controller) TestHttpServer {
	router := utils.GetDefaultRouter()
	controller.RegisterRoutes(router)
	server := httptest.NewServer(router)
	client := server.Client()
	return TestHttpServer{
		Server: server,
		client: client,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (c *TestHttpServer) GetURL(t *testing.T, path string) string {
	// The current implementation doesn't handle query parameters correctly
	// url.JoinPath will escape query parameters, treating them as part of the path

	// Split the path and query
	pathParts := strings.SplitN(path, "?", 2)
	basePath := pathParts[0]

	baseURL, err := url.JoinPath(c.Server.URL, basePath)
	if err != nil {
		t.Fatal(err)
	}

	// If there's a query string, append it
	if len(pathParts) > 1 {
		return baseURL + "?" + pathParts[1]
	}

	return baseURL
}

func (c *TestHttpServer) MustDo(
	t *testing.T,
	method string,
	url string,
	reqBody interface{},
	respBody interface{},
) int {
	// Encode request
	var buf bytes.Buffer
	if reqBody != nil {
		err := json.NewEncoder(&buf).Encode(reqBody)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Prepare request
	req, err := http.NewRequest(method, c.GetURL(t, url), &buf)
	if err != nil {
		t.Fatal(err)
	}

	// Do request
	resp, err := c.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the entire response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Handle response body based on its type
	if respBody != nil {
		// If respBody is a string pointer, assign the body as is
		if strPtr, ok := respBody.(*string); ok {
			*strPtr = string(bodyBytes)
		} else {
			// Otherwise try to parse as JSON
			err = json.Unmarshal(bodyBytes, respBody)
			if err != nil {
				t.Logf("resp.Body: %s", string(bodyBytes))
				t.Fatalf("status: %d, err:%s", resp.StatusCode, err)
			}
		}
	}

	return resp.StatusCode
}

////////////////////////////////////////////////////////////////////////////////

func (c *TestHttpServer) MustDoAndMatchCode(
	t *testing.T,
	method string,
	url string,
	reqBody any,
	respBody any,
	code int,
) {
	respCode := c.MustDo(t, method, url, reqBody, respBody)
	convey.So(respCode, convey.ShouldEqual, code)
}

/**
func (c *TestHttpServer) MustSucceeded(t *testing.T, method string, url string, reqBody interface{}, respBody interface{}, code int) {
	respCode := c.MustDo(t, method, url, reqBody, respBody)
	convey.So(respCode, convey.ShouldEqual, code)
}

func (c *TestHttpServer) MustFailed(t *testing.T, method string, path string, reqBody interface{}, err *errors.Error) {
	respBody := new(errors.Error)
	respCode := c.MustDo(t, method, path, reqBody, respBody)
	convey.So(respBody.Message, convey.ShouldEqual, err.Message)
	convey.So(respCode, convey.ShouldEqual, err.Code)
}
*/
