package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WangWilly/swap-estimation/pkgs/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIdContext(t *testing.T) {
	Convey("Given the context ID functions", t, func() {
		baseCtx := context.Background()

		Convey("When using ctxWithID to add an xid.ID to context", func() {
			id := xid.New()
			ctx := ctxWithID(baseCtx, id)

			Convey("Then IdFromCtx should retrieve the ID correctly", func() {
				retrievedID, ok := IdFromCtx(ctx)
				So(ok, ShouldBeTrue)
				So(retrievedID, ShouldEqual, id.String())
			})
		})

		Convey("When using ctxWithStringID to add a string ID to context", func() {
			idStr := "test-request-id"
			ctx := ctxWithStringID(baseCtx, idStr)

			Convey("Then IdFromCtx should retrieve the string ID correctly", func() {
				retrievedID, ok := IdFromCtx(ctx)
				So(ok, ShouldBeTrue)
				So(retrievedID, ShouldEqual, idStr)
			})
		})

		Convey("When the context has no ID", func() {
			Convey("Then IdFromCtx should return not ok", func() {
				retrievedID, ok := IdFromCtx(baseCtx)
				So(ok, ShouldBeFalse)
				So(retrievedID, ShouldEqual, "")
			})
		})
	})
}

func TestLoggingMiddleware(t *testing.T) {
	Convey("Given the LoggingMiddleware", t, func() {
		gin.SetMode(gin.TestMode)
		utils.InitLogging(context.Background())

		Convey("When a request has no request ID", func() {
			middleware := LoggingMiddleware()

			// Variables to capture values from handler
			var handlerHasID bool
			var handlerID string

			// Create a test request
			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Set up the gin context
			_, r := gin.CreateTestContext(w)
			r.Use(middleware)
			r.GET("/test", func(c *gin.Context) {
				// Get the ID from the context to verify it was set
				id, ok := IdFromCtx(c.Request.Context())

				// Store these values for testing
				handlerHasID = ok
				handlerID = id

				c.Status(200)
			})

			r.ServeHTTP(w, req)

			Convey("Then it should generate and set a request ID", func() {
				So(req.Header.Get(utils.RequestIdHeader), ShouldNotBeEmpty)
				So(handlerHasID, ShouldBeTrue)
				So(handlerID, ShouldNotBeEmpty)
				So(handlerID, ShouldEqual, req.Header.Get(utils.RequestIdHeader))
			})
		})

		Convey("When a request already has a request ID", func() {
			middleware := LoggingMiddleware()

			// Variables to capture values from handler
			var handlerHasID bool
			var handlerID string

			// Create a test request with ID
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set(utils.RequestIdHeader, "existing-id")
			w := httptest.NewRecorder()

			// Set up the gin context
			_, r := gin.CreateTestContext(w)
			r.Use(middleware)
			r.GET("/test", func(c *gin.Context) {
				// Get the ID from the context
				id, ok := IdFromCtx(c.Request.Context())

				// Store these values for testing
				handlerHasID = ok
				handlerID = id

				c.Status(200)
			})

			r.ServeHTTP(w, req)

			Convey("Then it should use the existing request ID", func() {
				So(req.Header.Get(utils.RequestIdHeader), ShouldEqual, "existing-id")
				So(handlerHasID, ShouldBeTrue)
				So(handlerID, ShouldEqual, "existing-id")
			})
		})
	})
}
