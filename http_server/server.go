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
	"fmt"
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
	s.backend = make(map[string]int)
	go s.probe_backends(probe)
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h := s.handler(r); h != nil {
		client := xff(r)
		hpr_utils.Log(fmt.Sprintf("Request from: %s Url: %s", client, r.Host))
		r.Header.Add("X-Forwarded-For‎", client)
		r.Header.Add("X-Real-IP", client)
		h.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Not found.", http.StatusNotFound)
}

func (s *Server) handler(req *http.Request) http.Handler {
	h := req.Host
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}

	_, ok := s.proxy[h]
	if !ok {
		f, _ := rc.ZRange(fmt.Sprintf("hpr-backends::%s", h), 0, -1, true)
		if len(f) == 0 {
			return nil
		}

		var url string
		for _, be := range f {
			count, err := strconv.Atoi(be)
			if err != nil {
				url = be
				continue
			}

			for r := 0; r <= count; r++ {
				s.mu.Lock()
				s.proxy[h] = append(s.proxy[h], Proxy{0, url, makeHandler(url)})
				s.mu.Unlock()
			}
		}

	}
	return s.Next(h)
}

/* TODO: Implement more balance algorithms */
func (s *Server) Next(h string) http.Handler {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.backend[h]++
	total := len(s.proxy[h])
	if s.backend[h] == total {
		s.backend[h] = 0
	}
	hpr_utils.Log(fmt.Sprintf("Using backend: %s Url: %s", s.proxy[h][s.backend[h]].Backend, h))
	return s.proxy[h][s.backend[h]].handler
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func (s *Server) probe_backends(probe time.Duration) {
	timeout := time.Duration(2 * time.Second)
	transport := http.Transport{Dial: dialTimeout}
	client := &http.Client{
		Transport: &transport,
	}

	for {
		time.Sleep(probe)
		// s.mu.Lock()
		for vhost, backends := range s.proxy {
			for backend := range backends {
				hpr_utils.Log(fmt.Sprintf("vhost: %s backend: %s", vhost, backend.Backend))
			}
			//client.Get(s.proxy[k][s.backend[h]])
		}
		// s.mu.Unlock()
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
