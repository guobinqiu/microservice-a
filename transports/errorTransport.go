package transports

import (
	"context"
	"encoding/json"
	"microservice-a/models"
	"microservice-a/myerror"
	"net/http"
)

func ErrorEncoder(c context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
	if err, ok := err.(*myerror.Error); ok {
		w.WriteHeader(err.StatusCode())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(models.Error{Error: err.Error()})
}
