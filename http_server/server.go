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
	"errors"
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
	proxy   map[string][]Proxy
	backend map[string]int
}

type Proxy struct {
	//	last    time.Time
	Backend string
	handler http.Handler
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
	vhost := req.Host
	if i := strings.Index(vhost, ":"); i >= 0 {
		vhost = vhost[:i]
	}

	s.mu.RLock()
	_, ok := s.proxy[vhost]
	if !ok {
		err := s.populate_proxies(vhost)
		if err != nil {
			hpr_utils.Log(fmt.Sprintf("%s for vhost %s", err, vhost))
			return nil
		}
	}
	s.mu.RUnlock()
	return s.Next(vhost)
}

func (s *Server) populate_proxies(vhost string) (err error) {
	f, _ := rc.ZRange(fmt.Sprintf("hpr-backends::%s", vhost), 0, -1, true)
	if len(f) == 0 {
		return errors.New("Backend list is empty")
	}

	var url string
	for _, be := range f {
		count, err := strconv.Atoi(be)
		if err != nil {
			url = be
			continue
		}

		for r := 1; r <= count; r++ {
			s.proxy[vhost] = append(s.proxy[vhost],
				Proxy{fmt.Sprintf("http://%s", url), makeHandler(url)})
		}
	}
	return
}

/* TODO: Implement more balance algorithms */
func (s *Server) Next(vhost string) http.Handler {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.backend[vhost]++
	total := len(s.proxy[vhost])
	if s.backend[vhost] >= total {
		s.backend[vhost] = 0
	}
	hpr_utils.Log(fmt.Sprintf(
		"Using backend: %s Url: %s", s.proxy[vhost][s.backend[vhost]].Backend, vhost))
	return s.proxy[vhost][s.backend[vhost]].handler
}

/* TODO: Implement more probes */
func (s *Server) probe_backends(probe time.Duration) {
	transport := http.Transport{Dial: dialTimeout}
	client := &http.Client{
		Transport: &transport,
	}

	for {
		time.Sleep(probe)
		s.mu.Lock()
		for vhost, backends := range s.proxy {
			// err := s.populate_proxies(vhost)
			is_dead := make(map[string]bool)
			removed := 0
			for backend := range backends {
				backend = backend - removed
				hpr_utils.Log(fmt.Sprintf(
					"vhost: %s backends: %s", vhost, s.proxy[vhost][backend].Backend))
				if is_dead[s.proxy[vhost][backend].Backend] == true {
					hpr_utils.Log(fmt.Sprintf("Removing dead backend: %s", s.proxy[vhost][backend].Backend))
					s.proxy[vhost] = s.proxy[vhost][:backend+copy(s.proxy[vhost][backend:], s.proxy[vhost][backend+1:])]
					removed++
					continue
				}

				_, err := client.Get(s.proxy[vhost][backend].Backend)
				if err != nil {
					hpr_utils.Log(fmt.Sprintf("Removing dead backend: %s with error %s", s.proxy[vhost][backend].Backend, err))
					is_dead[s.proxy[vhost][backend].Backend] = true
					s.proxy[vhost] = s.proxy[vhost][:backend+copy(s.proxy[vhost][backend:], s.proxy[vhost][backend+1:])]
					removed++
				} else {
					hpr_utils.Log(fmt.Sprintf("Alive: %s", s.proxy[vhost][backend].Backend))

				}
			}

		}
		s.mu.Unlock()
	}
}

func dialTimeout(network, addr string) (net.Conn, error) {
	timeout := time.Duration(2 * time.Second)
	return net.DialTimeout(network, addr, timeout)
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
