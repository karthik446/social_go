package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal-server-error path: %s, method:%s,  %s", r.URL.Path, r.Method, err)
	writeJSONError(w, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad-request-error path: %s, method:%s,  %s", r.URL.Path, r.Method, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not-found-error path: %s, method:%s,  %s", r.URL.Path, r.Method, err)
	writeJSONError(w, http.StatusNotFound, err.Error())
}

func (app *application) duplicateKeyConflict(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("duplicate-key-conflict path: %s, method:%s,  %s", r.URL.Path, r.Method, err)
	writeJSONError(w, http.StatusConflict, err.Error())
}
