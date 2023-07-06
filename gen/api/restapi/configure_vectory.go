// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
	"strings"

	"Vectory/gen/api/restapi/operations"
	"Vectory/gen/api/restapi/operations/collection"
)

//go:generate swagger generate server --target ../../api --name Vectory --spec ../../../../../../../var/folders/lp/s7st0l3n27v75m180d9wbtcm0000gn/T/spec.yaml133420447 --principal string --exclude-main

func configureFlags(api *operations.VectoryAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.VectoryAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.CollectionAddCollectionHandler == nil {
		api.CollectionAddCollectionHandler = collection.AddCollectionHandlerFunc(func(params collection.AddCollectionParams) middleware.Responder {
			return middleware.NotImplemented("operation collection.AddCollection has not yet been implemented")
		})
	}
	if api.CollectionDeleteCollectionHandler == nil {
		api.CollectionDeleteCollectionHandler = collection.DeleteCollectionHandlerFunc(func(params collection.DeleteCollectionParams) middleware.Responder {
			return middleware.NotImplemented("operation collection.DeleteCollection has not yet been implemented")
		})
	}
	if api.CollectionGetCollectionHandler == nil {
		api.CollectionGetCollectionHandler = collection.GetCollectionHandlerFunc(func(params collection.GetCollectionParams) middleware.Responder {
			return middleware.NotImplemented("operation collection.GetCollection has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Shortcut helpers for swagger-ui
		switch r.URL.Path {
		case "/swagger-ui", "/api/help", "/doc", "/doc/":
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("docs/swagger-ui"))).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
