package bait

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

type BaitServer struct {
	config map[string]*BaitConfig
}

func NewServer(ctx context.Context, config []*BaitConfig) *http.Server {
	mapped := map[string]*BaitConfig{}
	for _, c := range config {
		mapped[c.Request] = c
	}
	s := &BaitServer{
		config: mapped,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)
	srv := &http.Server{
		Addr:    ":2248",
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	return srv
}

func (s *BaitServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		dur := time.Since(start)
		w.Header().Set("X-Response-Time", dur.String())
		log.Printf("%-8s %s", r.URL.Path, dur)
	}()

	c, ok := s.config[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := c.Cmd(ctx)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	defer func() {
		w.Write([]byte("\n>>> stdout\n"))
		w.Write(stdout.Bytes())
		w.Write([]byte("\n>>> stderr\n"))
		w.Write(stderr.Bytes())
		w.Write([]byte("\n"))
	}()

	c.Lock()
	err := cmd.Run()
	c.Unlock()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
