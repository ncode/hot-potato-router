# HPR - Hot Potato Router
###### dynamic http reverse proxy made easy 
## How it works?

HPR receives a connection from interwebs, looks for backends on Redis for a given vhost and proxy the connection among all backends.

<img src="https://raw.github.com/ncode/hot-potato-router/master/hpr.png">

## Config:
### hpr.conf

    hpr:
      http_addr: :80                 # bind address
      https_addr:                    # https_addres
      cert_file:                     # cert file for https
      key_file:                      # key file for https
      probe_interval: 50             # seconds between probes
      db_backend: redis              # backend type

    redis:
      server_list: 127.0.0.1:6379    # list of redis servers to connect

## Instalation:
### Building

    $ git clone https://github.com/ncode/hot-potato-router.git
    $ cd hot-potato-router
    $ go get -v .
    $ make build

### Packaging on Debian and Ubuntu

    $ apt-get install golang
    $ git clone https://github.com/ncode/hot-potato-router.git
    $ cd gogix hot-potato-router
    $ dpkg-buildpackage -us -uc -rfakeroot

## Cli:
### hprctl

    $ hprctl -h
    Hot Potato Router Control

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

    $ hprctl add hpr.martinez.io 127.0.0.1:8080
    $ hprctl add hpr.martinez.io 127.0.0.1:8081
    $ hprctl show hpr.martinez.io
    :: vhost [ hpr.martinez.io ]
    -- backend 127.0.0.1:8080 weight=1
    -- backend 127.0.0.1:8081 weight=1

    $ hprctl list
    :: vhost [ hpr.martinez.io ]

## Depends:
* go-redis - https://github.com/fiorix/go-redis.git
* gdocopt  - https://github.com/docopt/docopt.go
