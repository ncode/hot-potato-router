package main

import (
	"fmt"
	docopt "github.com/docopt/docopt.go"
	"github.com/fiorix/go-redis/redis"
	hpr_utils "github.com/ncode/hot-potato-router/utils"
	"net/http"
	"strconv"
	"time"
)

var (
	cfg = hpr_utils.NewConfig()
)

func main() {
	usage := `Hot Potato Router Control.

Usage:
  hprctl add <name> backend <ip:port> [--weight=<n>]
  hprctl del <name> backend <ip:port> [--weight=<n>]
  hprctl show <name> 
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
  --drifting    Drifting mine.`

	arguments, _ := docopt.Parse(usage, nil, true, "Hot Potato Router 0.3.0", false)
	fmt.Println(arguments)
}
