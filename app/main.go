package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/IvanSkripnikov/golang_otus_project/controllers"
	"github.com/IvanSkripnikov/golang_otus_project/logger"
)

type route struct {
	method string

	regex *regexp.Regexp

	handler http.HandlerFunc
}

var routes = []route{
	newRoute(http.MethodGet, "/", controllers.HelloPageHandler),

	newRoute(http.MethodGet, "/banners", controllers.BannersHandler),

	newRoute(http.MethodGet, "/banners/([0-9]+)", controllers.BannerHandler),

	newRoute(http.MethodPost, "/add_banner_to_slot/([\\S]+)", controllers.AddBannerHandler),

	newRoute(http.MethodDelete, "/remove_banner_from_slot/([\\S]+)", controllers.RemoveBannerHandler),

	newRoute(http.MethodGet, "/get_banner_for_show/([\\S]+)", controllers.GetBannerForShowHandler),

	newRoute(http.MethodPost, "/event_click/([\\S]+)", controllers.ClickHandler),
}

func initHTTPServer() {
	http.HandleFunc("/", Serve)

	err := http.ListenAndServe(":8080", nil) //nolint:gosec
	if err != nil {
		errMessage := fmt.Sprintf("Can't init HTTP server: %v", err)

		logger.SendToErrorLog(errMessage)
	}
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string

	found := false

	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)

				continue
			}

			found = true

			route.handler(w, r)
		}
	}

	if !found && len(allow) == 0 {
		w.WriteHeader(http.StatusNotFound)

		http.NotFound(w, r)

		return
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))

		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)

		return
	}
}

func GetHTTPHandler() *http.ServeMux {
	httpHandler := http.NewServeMux()

	for _, route := range routes {
		httpHandler.HandleFunc(handleRegexp(route.regex), route.handler)
	}

	return httpHandler
}

func handleRegexp(regExp *regexp.Regexp) string {
	expr := regExp.String()[1 : len(regExp.String())-1]

	var result string

	if strings.Count(expr, "/") > 1 {
		parts := strings.Split(expr, "/")

		parts = parts[:len(parts)-1]

		result = strings.Join(parts, "/") + "/"
	} else {
		result = expr
	}

	return result
}

func main() {
	initHTTPServer()
}
