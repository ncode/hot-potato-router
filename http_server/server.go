/*
   Copyright 2013 Juliano Martinez <juliano@martinez.io>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

      Based on http://github.com/nf/webfront

   @author: Juliano Martinez
*/

package http_server

import (
	"github.com/fiorix/go-redis/redis"
	hpr_utils "github.com/ncode/hot-potato-router/utils"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	cfg = hpr_utils.NewConfig()
	rc  = redis.New(cfg.Options["redis"]["server_list"])
)

func xff(req *http.Request) string {
	remote_addr := strings.Split(req.RemoteAddr, ":")
	if len(remote_addr) == 0 {
		return ""
	}
	return remote_addr[0]
}

type Server struct {
	mu      sync.RWMutex
	last    time.Time
	proxy   map[string][]Proxy
	backend map[string]int
}

type Proxy struct {
	Connections int64
	Backend     string
	handler     http.Handler
}

func Listen(fd int, addr string) net.Listener {
	var l net.Listener
	var err error
	if fd >= 3 {
		l, err = net.FileListener(os.NewFile(uintptr(fd), "http"))
	} else {
		l, err = net.Listen("tcp", addr)
	}
	if err != nil {
		log.Fatal(err)
	}
	return l
}

func NewServer(probe time.Duration) (*Server, error) {
	s := new(Server)
	s.proxy = make(map[string][]Proxy)
	go s.probe_backends(probe)
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h := s.handler(r); h != nil {
		r.Header.Add("X-Forwarded-For‎", xff(r))
		r.Header.Add("X-Real-IP", xff(r))
		h.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not found.", http.StatusNotFound)
}

func (s *Server) handler(req *http.Request) http.Handler {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h := req.Host
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}

	_, ok := s.proxy[h]
	if !ok {
		f, _ := rc.Get(h)
		if f == "" {
			return nil
		}
		s.proxy[h] = append(s.proxy[h], Proxy{0, f, makeHandler(f)})
	}
	return s.Next(h)
}

func (s *Server) Next(h string) http.Handler {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.backend[h]++
	total := len(s.proxy[h])
	if s.backend[h] == total {
		s.backend[h] = 0
	}
	return s.proxy[h][s.backend[h]].handler
}

func (s *Server) probe_backends(probe time.Duration) {
	for {
		s.mu.Lock()
		s.mu.Unlock()
		time.Sleep(probe)
	}
}

func makeHandler(f string) http.Handler {
	if f != "" {
		return &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = f
			},
		}
	}
	return nil
}