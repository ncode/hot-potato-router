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
	"fmt"
	docopt "github.com/docopt/docopt.go"
	"github.com/fiorix/go-redis/redis"
	hpr_utils "github.com/ncode/hot-potato-router/utils"
	"strconv"
)

var (
	cfg = hpr_utils.NewConfig()
	rc  = redis.New(cfg.Options["redis"]["server_list"])
)

func main() {
	usage := `Hot Potato Router Control.

Usage:
  hprctl add <vhost> <backend_ip:port> [--weight=<n>]
  hprctl del <vhost> <backend_ip:port> [--weight=<n>]
  hprctl show <vhost> 
  hprctl list

Args:
  add           add a new vhost and backend
  dell          del a vhost and a backend
  show          show all backends from a given vhost
  list          list all vhosts

Options:
  -h --help     Show this screen.
  --version     Show version.
  --weight=<n>  Weight in wrr [default: 1].
`

	arguments, _ := docopt.Parse(usage, nil, true, "Hot Potato Router 0.3.0", false)
	if arguments["add"] == true {
		_, err := rc.ZAdd(fmt.Sprintf("hpr-backends::%s", arguments["<vhost>"]), arguments["--weight"], arguments["<backend_ip:port>"])
		hpr_utils.CheckPanic(err, "Unable to write on hpr database")
		return
	}

	/* if arguments["del"] == true {
		_, err := rc.ZRem(fmt.Sprintf("hpr-backends::%s", arguments["<vhost>"]), arguments["--weight"], arguments["<backend_ip:port>"])
		hpr_utils.CheckPanic(err, "Unable to write on hpr database")
		return
	} */

	if arguments["show"] == true {
		bes, err := rc.ZRange(fmt.Sprintf("hpr-backends::%s", arguments["<vhost>"]), 0, -1, true)
		hpr_utils.CheckPanic(err, "Unable to write on hpr database")
		var url string
		for _, be := range bes {
			count, err := strconv.Atoi(be)
			if err != nil {
				url = be
				continue
			}
			fmt.Printf("backend %s wight %s", url, count)
		}
		return
	}

	if arguments["list"] == true {
		keys, err := rc.Keys("hpr-backends::*")
		hpr_utils.CheckPanic(err, "Unable to write on hpr database")
		for _, k := range keys {
			fmt.Printf("vhost %s", k)
		}
		return
	}
}
