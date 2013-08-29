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

   @author: Juliano Martinez
*/

package main

import (
	"crypto/tls"
	"github.com/fiorix/go-redis/redis"
	hpr_http_server "github.com/ncode/hot-potato-router/http_server"
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
)

func main() {
	probe_interval, _ := strconv.Atoi(cfg.Options["hpr"]["probe_interval"])
	if probe_interval == 0 {
		probe_interval = 10
	}
	s, err := hpr_http_server.NewServer(time.Duration(probe_interval) * time.Second)
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
		l := tls.NewListener(hpr_http_server.Listen(https_fd, cfg.Options["hpr"]["https_addr"]), c)
		go func() {
			log.Fatal(http.Serve(l, s))
		}()
	}
	log.Fatal(http.Serve(hpr_http_server.Listen(http_fd, cfg.Options["hpr"]["http_addr"]), s))
}
