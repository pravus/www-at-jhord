package main

import (
  "context"
  "errors"
  "log"
  "net"
  "os"

  "google.golang.org/grpc"
  "github.com/pravus/www-at-jhord/registry"
)

var MAX_VISITS = 10

type Entry struct {
  count uint64
  log   []*registry.Visit
}

type VisitRegistry struct {
  register map [string] *Entry
}

func main() {
  bind := _env("GRPC_BIND", ":5000")
  tcp, err := net.Listen("tcp", bind)
  if err != nil {
    log.Fatalf("grpc.listen error=%v", err)
  }

  server := grpc.NewServer()
  registry.RegisterRegistryServer(server, &VisitRegistry{
    register: map [string] *Entry{},
  })

  log.Printf("grpc.serve bind=%s", bind)
  server.Serve(tcp)
}

func _env(name, ifEmpty string) string {
  value := os.Getenv(name)
  if value == "" {
    value = ifEmpty
  }
  return value
}

func (vr *VisitRegistry) GetMaxVisits(ctx context.Context, req *registry.GetMaxVisitsRequest) (*registry.GetMaxVisitsResponse, error) {
  max := uint64(MAX_VISITS)
  return &registry.GetMaxVisitsResponse{Max: &max}, nil
}

func (vr *VisitRegistry) LogVisit(ctx context.Context, visit *registry.Visit) (*registry.LogVisitResponse, error) {
  if visit == nil {
    return nil, errors.New("No visit")
  }
  if visit.Address == nil || *visit.Address == "" {
    return nil, errors.New("No address")
  }
  entry, found := vr.register[*visit.Address]
  if !found {
    entry = &Entry{
      count: 0,
      log:   []*registry.Visit{},
    }
    vr.register[*visit.Address] = entry
  }
  entry.count += 1
  entry.log = append(entry.log, visit)
  cut := len(entry.log) - MAX_VISITS
  if cut > 0 {
    entry.log = entry.log[cut:]
  }
  log.Printf("visit %s %s %d", *visit.Address, *visit.Agent, entry.count)
  return &registry.LogVisitResponse{Count: &entry.count}, nil
}

func (vr *VisitRegistry) GetVisits(req *registry.GetVisitsRequest, stream registry.Registry_GetVisitsServer) error {
  if req == nil {
    return errors.New("No request")
  }
  if req.Address == nil || *req.Address == "" {
    return errors.New("No address")
  }
  entry, found := vr.register[*req.Address]
  if !found {
    return errors.New("Address not found")
  }
  for _, visit := range entry.log {
    err := stream.Send(visit)
    if err != nil {
      return err
    }
  }
  return nil
}
