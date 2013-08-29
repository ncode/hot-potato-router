# HPR - Hot Potato Router
###### dynamic http reverse proxy made easy 
## How it works?

HPR receives a connection from interwebs look for the destination on Redis and proxy the connection for one or more instances, to avoid a dead instance serving content we have a goroutine checking and updating the instance state on redis.

<img src="https://raw.github.com/ncode/hot-potato-router/master/hpr.png">
