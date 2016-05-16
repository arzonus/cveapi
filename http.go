package main

import (
	"github.com/arzonus/cveapi/interfaces"
	"github.com/gorilla/mux"
	"net/http"
)

func HTTP(wsh interfaces.WebserviceHandler) {
	r := mux.NewRouter()
	r.HandleFunc("/cve/{cveId}", route(wsh.GetCVEById)).Methods("GET")
	r.HandleFunc("/product/{name}/{version}", route(wsh.GetListCVEByProduct)).Methods("GET")
	http.ListenAndServe(":3000", r)
}

func route(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {

	headers := func(h http.HandlerFunc) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			h.ServeHTTP(res, req)
		}
	}

	h = headers(h)
	for _, m := range middleware {
		h = m(h)
	}

	return h
}
