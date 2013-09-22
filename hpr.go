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
	"github.com/ncode/hot-potato-router/http_server"
	"github.com/ncode/hot-potato-router/utils"
	"net/http"
	"strconv"
	"time"
)

var (
	cfg = utils.NewConfig()
)

func main() {
	utils.Log("Starting Hot Potato Router...")
	probe_interval, _ := strconv.Atoi(cfg.Options["hpr"]["probe_interval"])
	if probe_interval == 0 {
		probe_interval = 10
	}
	s, err := http_server.NewServer(time.Duration(probe_interval) * time.Second)
	utils.CheckPanic(err, "Unable to spawn")

	if cfg.Options["hpr"]["https_addr"] != "" {
		cert, err := tls.LoadX509KeyPair(cfg.Options["hpr"]["cert_file"], cfg.Options["hpr"]["key_file"])
		utils.CheckPanic(err, "Unable to load certificate")

		c := &tls.Config{Certificates: []tls.Certificate{cert}}
		l := tls.NewListener(http_server.Listen(fg.Options["hpr"]["https_addr"]), c)
		go func() {
			utils.CheckPanic(http.Serve(l, s), "Problem with https server")
		}()
	}
	utils.CheckPanic(
		http.Serve(http_server.Listen(cfg.Options["hpr"]["http_addr"]), s),
		"Problem with http server")
}
