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
	//	"github.com/fiorix/go-redis/redis"
	hpr_utils "github.com/ncode/hot-potato-router/utils"
)

var (
	cfg = hpr_utils.NewConfig()
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
  list          list all backends and vhosts

Options:
  -h --help     Show this screen.
  --version     Show version.
  --weight=<n>  Weight in wrr [default: 1].
`

	arguments, _ := docopt.Parse(usage, nil, true, "Hot Potato Router 0.3.0", false)
	fmt.Println(arguments)
}
