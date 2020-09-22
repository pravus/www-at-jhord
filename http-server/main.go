package main

import (
  "log"
  "net/http"
  "os"

  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
)

var GRPC_DIAL = "www-at-jhord-grpc:5000"

func main() {
  bind := _env("HTTP_BIND", ":8000")
  mux  := chi.NewRouter()

  mux.Use(middleware.Logger)
  mux.Get("/healthz", http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {w.Write([]byte("healthy\n"));}));
  mux.Route("/",       RootRouter())
  mux.Route("/visits", VisitsRouter())
  mux.Route("/resume", ResumeRouter())

  log.Printf("http up bind=%s", bind)
  err := http.ListenAndServe(bind, mux)
  if err != nil {
    log.Fatalf("http error=%v", err)
  }
}

func _env(name, ifEmpty string) string {
  value := os.Getenv(name)
  if value == "" {
    value = ifEmpty
  }
  return value
}
