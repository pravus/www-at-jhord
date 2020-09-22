package main

import (
  "context"
  "html/template"
  "io"
  "log"
  "net/http"
  "time"
  "strings"

  "google.golang.org/grpc"
  "github.com/go-chi/chi"

  "github.com/pravus/www-at-jhord/registry"
)

type Visit struct {
  Date  string
  Agent string
}

type VisitsInfo struct {
  Max     uint64
  Address string
  Count   uint64
  Visits  []*Visit
}

func VisitsRouter() func (chi.Router) {
  return func (r chi.Router) {
    r.Get("/", Visits())
  }
}

func Visits() http.HandlerFunc {
  content := template.Must(template.ParseFiles("visits/content.html"))
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/html")

    address := r.RemoteAddr
    index   := strings.Index(address, ":")
    if index >= 0 {
      address = address[0:index]
    }

    err := content.Execute(w, getVisitsInfo(address, r.Header.Get("User-Agent")))
    if err != nil {
      log.Printf("template: error=%v", err)
      return
    }
  })
}

func getVisitsInfo(address, agent string) *VisitsInfo {
  info := &VisitsInfo{
    Address: address,
    Count:   0,
    Visits:  nil,
  }
  tcp, err := grpc.Dial(GRPC_DIAL, grpc.WithInsecure())
  if err != nil {
    log.Printf("grpc.dial error=%v", err)
    return info
  }
  defer tcp.Close()

  ts := time.Now().UnixNano()
  visit := &registry.Visit{
    Address: &address,
    Agent:   &agent,
    Ts:      &ts,
  }
  client := registry.NewRegistryClient(tcp)

  limit, err := client.GetMaxVisits(context.Background(), &registry.GetMaxVisitsRequest{})
  if err != nil {
    log.Printf("grpc.GetMaxVisits error=%v", err)
    return info
  }

  res, err := client.LogVisit(context.Background(), visit)
  if err != nil {
    log.Printf("grpc.LogVisit error=%v", err)
    return info
  }

  stream, err := client.GetVisits(context.Background(), &registry.GetVisitsRequest{Address: &address})
  if err != nil {
    log.Printf("grpc.GetVisits error=%v", err)
    return info
  }
  visits := []*Visit{}
  for {
    visit, err := stream.Recv()
    if err == io.EOF {
      break
    }
    if err != nil {
      log.Printf("grpc.GetVisits.Recv error=%v", err)
      break
    }
    visits = append(visits, &Visit{
      Date:  time.Unix(0, *visit.Ts).Format("2006-01-02 15:04:05"),
      Agent: *visit.Agent,
    })
  }

  info.Max    = *limit.Max
  info.Count  = *res.Count
  info.Visits = visits
  return info
}
