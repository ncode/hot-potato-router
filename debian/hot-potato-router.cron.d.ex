#
# Regular cron jobs for the hot-potato-router package
#
0 4	* * *	root	[ -x /usr/bin/hot-potato-router_maintenance ] && /usr/bin/hot-potato-router_maintenance
