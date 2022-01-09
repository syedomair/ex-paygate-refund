package routes

import (
	//"context"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/syedomair/ex-paygate-lib/lib/container"
	log "github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-refund/routes/refund"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// NewRouter comment
func NewRouter(c container.Container) *chi.Mux {
	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST, GET"},
		AllowedHeaders:   []string{"ApiKey", "RefreshToken", "Token", "FrontendURL", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/", routerSetup(
			c.Db(),
			c.Logger(),
			c.SigningKey(),
		))
	})

	return router
}

// Route Public
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Access      string
}

// Routes List Public
type Routes []Route

func routerSetup(db *gorm.DB, logger log.Logger, signingKey string) *chi.Mux {

	repoRefund := refund.NewPostgresRepository(db, logger)
	payRefund := refund.NewPaymentService(logger)

	router := chi.NewRouter()
	refundController := &refund.Controller{
		Logger: logger,
		Repo:   repoRefund,
		Pay:    payRefund,
	}

	var routes = Routes{
		Route{
			"Ping",
			"GET",
			"/ping",
			refundController.Ping,
			"public",
		},
		Route{
			"Refund",
			"POST",
			"/refund",
			refundController.RefundAction,
			"public",
		},
	}

	for _, route := range routes {

		handler := route.HandlerFunc

		/*
			if route.Access == "admin" {
				handler = securityWorkflowAdminMiddleware(strKey, logger, handler, signingKey)
			} else if route.Access == "networkadmin" {
				handler = securityNetworkAdminMiddleware(strKey, logger, handler, signingKey)
			} else if route.Access == "networkrobotadmin" {
				handler = securityNetworkRobotAdminMiddleware(strKey, logger, handler, signingKey)
			} else if route.Access == "user" {
				handler = securityUserMiddleware(strKey, logger, handler, signingKey)
			} else if route.Access == "globaladmin" {
				handler = securityGlobalAdminMiddleware(strKey, logger, handler, signingKey)
			}
		*/
		router.Method(route.Method, route.Pattern, handler)
	}
	return router
}
