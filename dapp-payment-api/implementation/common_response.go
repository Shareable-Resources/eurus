package implementation

import (
	"net/http"

	"github.com/go-chi/render"
)

func OK(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	status(w, r, 200, v...)
}

func Created(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	status(w, r, 201, v...)
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	status(w, r, 204)
}

func BadRequest(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	status(w, r, 400, v...)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	status(w, r, 401, v...)
}

func Forbidden(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	status(w, r, 403, v...)
}

func NotFound(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	status(w, r, 404, v...)
}

func Conflict(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	status(w, r, 409, v...)
}

func Gone(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	status(w, r, 410, v...)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, v ...interface{}) {
	status(w, r, 500, v...)
}

func status(w http.ResponseWriter, r *http.Request, status int, v ...interface{}) {
	// render.Respond() forces to have response body even if no payload is set
	// In this case just simply use the traditional way to return status code
	if len(v) > 0 {
		render.Status(r, status)
		render.Respond(w, r, v[0])
	} else {
		w.WriteHeader(status)
	}
}
