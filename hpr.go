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

package main

import (
	"crypto/tls"
	"github.com/fiorix/go-redis/redis"
	hpr_utils "github.com/ncode/hot-potato-router/utils"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	mu    sync.RWMutex
	last  time.Time
	proxy map[string]http.Handler
}

type Proxy struct {
	Backend string
	handler http.Handler
}

var (
	cfg = hpr_utils.NewConfig()
	rc  = redis.New(cfg.Options["redis"]["server_list"])
)

func main() {
	probe_interval, _ := strconv.Atoi(cfg.Options["hpr"]["probe_interval"])
	if probe_interval == 0 {
		probe_interval = 10
	}
	s, err := NewServer(time.Duration(probe_interval) * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	http_fd, _ := strconv.Atoi(cfg.Options["hpr"]["http_fd"])
	https_fd, _ := strconv.Atoi(cfg.Options["hpr"]["https_fd"])
	if https_fd >= 3 || cfg.Options["hpr"]["https_addr"] != "" {
		cert, err := tls.LoadX509KeyPair(cfg.Options["hpr"]["cert_file"], cfg.Options["hpr"]["key_file"])
		if err != nil {
			log.Fatal(err)
		}
		c := &tls.Config{Certificates: []tls.Certificate{cert}}
		l := tls.NewListener(listen(https_fd, cfg.Options["hpr"]["https_addr"]), c)
		go func() {
			log.Fatal(http.Serve(l, s))
		}()
	}
	log.Fatal(http.Serve(listen(http_fd, cfg.Options["hpr"]["http_addr"]), s))
}

func listen(fd int, addr string) net.Listener {
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
	s.proxy = make(map[string]http.Handler)
	go s.probe_backends(probe)
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h := s.handler(r); h != nil {
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
		v, _ := rc.Get(h)
		s.proxy[h] = makeHandler(v)
	}
	return s.proxy[h]
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
