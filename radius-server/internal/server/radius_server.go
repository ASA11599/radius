package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ASA11599/radius-server/internal/model"
	"github.com/ASA11599/radius-server/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type RadiusServer struct {
	store store.Store
	server *http.Server
}

func NewRadiusServer(host string, port int, s store.Store) *RadiusServer {
	rs := &RadiusServer{
		store: s,
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", host, port),
		},
	}
	rs.registerRouter()
	return rs
}

func (rs *RadiusServer) registerRouter() {

	apiRouter := chi.NewRouter()

	apiRouter.Use(
		middleware.CleanPath,
		middleware.Logger,
		middleware.RedirectSlashes,
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Content-Type"},
			ExposedHeaders: []string{"Link"},
			AllowCredentials: false,
			MaxAge: 300,
		}),
		middleware.AllowContentType("application/json"),
		setContentTypeJSON,
	)

	apiRouter.Get("/api/health", rs.healthCheck)
	apiRouter.Post("/api/posts", rs.createPost)
	apiRouter.Get("/api/posts", rs.getNearbyPosts)

	apiRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.ErrorResponse{ Message: "Not found" })
	})

	apiRouter.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ErrorResponse{ Message: "Method not allowed" })
	})

	mux := http.NewServeMux()

	mux.Handle("/api/", apiRouter)

	fileServer := http.FileServer(http.Dir("./dist/"))
	mux.Handle("/", http.StripPrefix("/", fileServer))

	rs.server.Handler = mux

}

func (rs *RadiusServer) Start() error {
	return rs.server.ListenAndServe()
}

func (rs *RadiusServer) Stop() error {
	return errors.Join(
		rs.server.Shutdown(context.Background()),
		rs.store.Close(),
	)
}

func (rs *RadiusServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{ "healthy": true })
}

func (rs *RadiusServer) createPost(w http.ResponseWriter, r *http.Request) {
	body, bodyError := io.ReadAll(r.Body)
	if bodyError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{ Message: bodyError.Error() })
		return
	}
	var postRequest model.PostRequest
	jsonError := json.Unmarshal(body, &postRequest)
	if jsonError != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{ Message: jsonError.Error() })
		return
	}
	if !postRequest.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{ Message: "Invalid request" })
		return
	}
	post := model.NewPost(postRequest)
	storeError := rs.store.SavePost(*post)
	if storeError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{ Message: "Storage error" })
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (rs *RadiusServer) getNearbyPosts(w http.ResponseWriter, r *http.Request) {
	lat, latError := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	long, longError := strconv.ParseFloat(r.URL.Query().Get("long"), 64)
	radius, radiusError := strconv.ParseFloat(r.URL.Query().Get("radius"), 64)
	if (latError != nil) || (longError != nil) || (radiusError != nil) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{ Message: "Error parsing query parameters" })
		return
	}
	loc := model.Location{ Latitude: lat, Longitude: long }
	if !loc.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{ Message: "Invalid location" })
		return
	}
	nearbyPosts, storeError := rs.store.GetNearbyPosts(loc, radius)
	if storeError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{ Message: "Storage error" })
		return
	}
	json.NewEncoder(w).Encode(nearbyPosts)
}
