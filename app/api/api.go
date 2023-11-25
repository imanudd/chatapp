package api

import (
	"context"
	"finalproject/config"
	mongoconn "finalproject/infra/db"
	"finalproject/internal/chatapp"
	customMiddleware "finalproject/middleware"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type API struct {
	cfg    *config.Config
	router *echo.Echo
	db     *mongo.Database
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func New() *API {
	cfg := config.LoadDefault()

	mongoDB := mongoconn.NewMongoConn(cfg)

	router := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("./assets/templates/*.html")),
	}
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Renderer = t

	api := &API{
		cfg:    cfg,
		router: router,
		db:     mongoDB,
	}
	return api
}

func (api API) BuildHandler() *echo.Echo {
	api.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	api.router.Use(customMiddleware.RequestIDContext())
	// api.router.HTTPErrorHandler = CustomHTTPErrorHandler(api.cfg, api.log, httplog.NewHTTPLog(api.db))

	chatapp.RegisterAPI(
		*api.router.Group(""),
		api.cfg,
		chatapp.NewService(api.cfg, api.db,
			chatapp.NewRepository(api.db, api.cfg)),
	)

	return api.router
}

func (api API) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	handleSigterm(func() {
		cancel()
	})

	// Start server
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", api.cfg.Server.PORT),
		Handler: api.BuildHandler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	fmt.Printf("server is running at port: %v [env: %v]", api.cfg.Server.PORT, api.cfg.Server.ENV)

	gracefulShutdownServer(ctx, &server)
}

func gracefulShutdownServer(ctx context.Context, srv *http.Server) {

	<-ctx.Done()

	fmt.Println("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown failed:%+s", err)
	}

	fmt.Println("server exited properly")

}

func handleSigterm(exitFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		exitFunc()
	}()
}
