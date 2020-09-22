package main

import (
  "context"
  "html/template"
  "log"
  "net/http"
  "time"
  "strings"

  "google.golang.org/grpc"
  "github.com/go-chi/chi"

  "github.com/pravus/www-at-jhord/registry"
)

func RootRouter() func (chi.Router) {
  return func (r chi.Router) {
    r.Get("/", Home())
  }
}

type RootVars struct {
  Prefix string
}

func Home() http.HandlerFunc {
  content := template.Must(template.ParseFiles("root/content.html"))
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    logVisit(r)
    w.Header().Add("Content-Type", "text/html")
    prefix := "@jhord/"
    if r.URL.Path[len(r.URL.Path) - 1:] == "/" {
      prefix = ""
    }
    err := content.Execute(w, &RootVars{Prefix: prefix})
    if err != nil {
      log.Printf("template: error=%v", err)
    }
  })
}

func logVisit(r *http.Request) {
  tcp, err := grpc.Dial(GRPC_DIAL, grpc.WithInsecure())
  if err != nil {
    log.Printf("grpc.dial error=%v", err)
    return
  }
  defer tcp.Close()

  address := r.RemoteAddr
  index   := strings.Index(address, ":")
  if index >= 0 {
    address = address[0:index]
  }
  agent  := r.Header.Get("User-Agent")
  ts     := time.Now().UnixNano()
  client := registry.NewRegistryClient(tcp)
  _, err = client.LogVisit(context.Background(), &registry.Visit{
    Address: &address,
    Agent:   &agent,
    Ts:      &ts,
  })
  if err != nil {
    log.Printf("grpc.LogVisit error=%v", err)
    return
  }
  return
}
