package routers

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCourseRouterRegistersOrderRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	CourseRouter(router.Group(""))

	routes := map[string]bool{}
	for _, route := range router.Routes() {
		routes[route.Method+" "+route.Path] = true
	}

	for _, route := range []string{
		"POST /course/:course_id/order",
		"POST /course/order/pay-callback",
		"GET /course/:course_id/access",
		"GET /course/:course_id/chapter/:chapter_id",
		"GET /course/:course_id/chapter/:chapter_id/:video_id",
		"GET /course/video/:video_id/hls-key",
	} {
		if !routes[route] {
			t.Fatalf("missing route %s", route)
		}
	}
}

func TestOrderRouterRegistersPayRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	OrderRouter(router.Group(""))

	routes := map[string]bool{}
	for _, route := range router.Routes() {
		routes[route.Method+" "+route.Path] = true
	}

	if !routes["POST /order/:order_no/pay"] {
		t.Fatalf("missing route POST /order/:order_no/pay")
	}
}
