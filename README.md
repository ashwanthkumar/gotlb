# GoTLB

GoTLB is a Go based TCP Load Balancer built for marathon apps.

## Why an another LB?
When you're doing micro-services there are number of load balancers available as choices like [Traefik](https://traefik.io/), [LinkerD](https://linkerd.io/), [HAProxy](https://www.haproxy.org/) via [marathon-lb](https://github.com/mesosphere/marathon-lb) or others, etc. But all of them support HTTP and some HTTP/2 and only one in that list TCP - HAProxy. Unfortunately it still doesn't hot reloading of the routes. I looked at things like [fabio](https://github.com/fabiolb/fabio) as well which recently added support for TCP. But that had an external dependency like Consul, which is something we don't have in our infrastructure. At Indix, we use application labels for configuring our apps. So the source of truth is always with the application's specification and not outside. Hence this is an attempt at solving these problems.

Also it was fun! :smile:

## Contribute
If you've any feature requests or issues, please open a Github issue. We accept PRs. Fork away!

## License
http://www.apache.org/licenses/LICENSE-2.0
