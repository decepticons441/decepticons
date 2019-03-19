package handlers

import (
	"log"
	"net/http"
)

/* TODO: implement a CORS middleware handler, as described
in https://drstearns.github.io/tutorials/cors/ that responds
with the following headers to all requests:

  Access-Control-Allow-Origin: *
  Access-Control-Allow-Methods: GET, PUT, POST, PATCH, DELETE
  Access-Control-Allow-Headers: Content-Type, Authorization
  Access-Control-Expose-Headers: Authorization
  Access-Control-Max-Age: 600
*/

// CorsHeader is a middleware handler that adds a header to the response
type CorsHandler struct {
	handler http.Handler
}

func (c *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Add("Access-Control-Expose-Headers", "Authorization")
	w.Header().Add("Access-Control-Max-Age", "600")

	log.Println("here")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
	} else {
		c.handler.ServeHTTP(w, r)
	}

	// if r.Method == http.MethodOptions {
	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	// 	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	// 	w.Header().Set("Access-Control-Max-Age", "600")
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// } else {
	// 	c.Handler.ServeHTTP(w, r)
	// }
}

func NewCorsHandler(wrapHandler http.Handler) *CorsHandler {
	return &CorsHandler{
		handler: wrapHandler,
	}
}
