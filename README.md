# HPR - Hot Potato Router
###### dynamic http reverse proxy made easy 
## How it works?

HPR receives a connection from interwebs look for the destination on Redis and proxy the connection for one or more instances, to avoid a dead instance serving content we have a goroutine checking and updating the instance state on redis.

<img src="https://raw.github.com/ncode/hot-potato-router/master/hpr.png">


## Usage:
###  hprctl

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


### hpr.conf

    hpr:
        http_addr: :80                 # bind address
        http_fd:
        https_addr:
        https_fd:
        cert_file:
        key_file:
        probe_interval: 50             # seconds between probes
        db_backend: redis              # backend type

    redis:
        server_list: 127.0.0.1:6379    # list of redis servers to connect

## Depends:
* go-redis - https://github.com/fiorix/go-redis.git
* gdocopt  - github.com/docopt/docopt.go
