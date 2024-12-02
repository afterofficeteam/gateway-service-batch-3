package routes

import (
	"gateway-service/config"
	"gateway-service/util/middleware"
	"log"
	"net/http"
	"strings"
	"time"

	cart "gateway-service/handlers/cart"
	order "gateway-service/handlers/order"
	user "gateway-service/handlers/users"

	"github.com/spf13/viper"
)

type Routes struct {
	Router *http.ServeMux
	User   *user.Handler
	Cart   *cart.Handler
	Order  *order.Handler
}

type RouteGroup struct {
	Router      *http.ServeMux
	Prefix      string
	Middlewares []func(http.Handler) http.Handler
}

func NewRouteGroup(router *http.ServeMux, prefix string) *RouteGroup {
	return &RouteGroup{
		Router:      router,
		Prefix:      prefix,
		Middlewares: []func(http.Handler) http.Handler{},
	}
}

func (rg *RouteGroup) Use(middlewares ...func(http.Handler) http.Handler) {
	rg.Middlewares = append(rg.Middlewares, middlewares...)
}

func (rg *RouteGroup) HandleFunc(methodAndPath string, handler http.HandlerFunc) {
	parts := strings.SplitN(methodAndPath, " ", 2)
	method := parts[0]
	path := rg.Prefix + parts[1]

	finalHandler := middleware.ApplyMiddleware(handler, rg.Middlewares...)

	rg.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.NotFound(w, r)
			return
		}
		finalHandler.ServeHTTP(w, r)
	})
}

func URLRewriter(baseURLPath string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, baseURLPath)
		next.ServeHTTP(w, r)
	}
}

func (r *Routes) setupBaseURL() {
	baseURL := viper.GetString("BASE_URL_PATH")
	if baseURL != "" && baseURL != "/" {
		r.Router.HandleFunc(baseURL+"/", URLRewriter(baseURL, r.Router))
	}
}

func (r *Routes) setupRouter() {
	r.Router = http.NewServeMux()
	r.setupBaseURL()

	userGroup := NewRouteGroup(r.Router, "/user")
	userGroup.Use(middleware.EnabledCors, middleware.LoggerMiddleware())
	userGroup.HandleFunc("POST /signup", r.User.SignUpByEmail)
	userGroup.HandleFunc("POST /signin", r.User.SignInByEmail)

	cartGroup := NewRouteGroup(r.Router, "/cart")
	cartGroup.Use(middleware.EnabledCors, middleware.LoggerMiddleware(), middleware.Authentication)
	cartGroup.HandleFunc("POST ", r.Cart.InsertCart)
	cartGroup.HandleFunc("GET /{id}", r.Cart.GetDetail)
	cartGroup.HandleFunc("DELETE /product/{product_id}", r.Cart.Delete)

	orderGroup := NewRouteGroup(r.Router, "/order")
	orderGroup.Use(middleware.EnabledCors, middleware.LoggerMiddleware(), middleware.Authentication)
	orderGroup.HandleFunc("POST ", r.Order.CreateOrder)
}

func (r *Routes) Run(port string) {
	r.setupRouter()

	log.Printf("[Running-Success] clients on localhost on port :%s", port)
	srv := &http.Server{
		Handler:      r.Router,
		Addr:         "localhost:" + port,
		WriteTimeout: config.WriteTimeout() * time.Second,
		ReadTimeout:  config.ReadTimeout() * time.Second,
	}

	log.Panic(srv.ListenAndServe())
}
