package main

import (
  "context"
  "log"
  "net/http"
  "os"
  "strings"

  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
)

var HTTP_BIND = ":8000"
var GRPC_DIAL = "www-at-jhord-grpc:5000"

func main() {
  bind := _env("HTTP_BIND", HTTP_BIND)
  mux  := chi.NewRouter()

  mux.Use(middleware.Logger)
  mux.Use(remoteAddr)

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

func remoteAddr(next http.Handler) http.Handler {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    remoteAddr := r.RemoteAddr
    forwardFor := r.Header.Get("X-Forwarded-For")
    if forwardFor != "" {
      remoteAddr = strings.Split(forwardFor, ", ")[0]
    }
    index := strings.LastIndex(remoteAddr, ":")
    if index >= 0 {
      remoteAddr = remoteAddr[0:index]
    }
    next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "RemoteAddr", remoteAddr)))
  })
}

func _env(name, ifEmpty string) string {
  value := os.Getenv(name)
  if value == "" {
    value = ifEmpty
  }
  return value
}
