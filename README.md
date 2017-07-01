# GoTLB

GoTLB is a Go based TCP Load Balancer built for marathon apps.

*Status*: **Work In Progress**

## Why an another LB?
When you're doing micro-services there are number of load balancers available as choices like [Traefik](https://traefik.io/), [LinkerD](https://linkerd.io/), [HAProxy](https://www.haproxy.org/) via [marathon-lb](https://github.com/mesosphere/marathon-lb) or others, etc. But all of them support HTTP and some HTTP/2 and only one in that list support TCP - HAProxy. Unfortunately it still doesn't hot reloading of the routes. We looked at things like [fabio](https://github.com/fabiolb/fabio) as well which recently added support for TCP. But that had an external dependency like Consul, which is something we don't have in our infrastructure. At Indix, we use application labels for configuring our apps. So the source of truth is always with the application's specification and not outside. Hence this is an attempt at solving these problems.

Also it was fun! :smile:

## Required Labels

We use the following list of labels in the app specification to identify and configure the frontends

| Property  | Description  |  Example  |
| :--- | :--- | :---: |
| tlb.enabled | Controls if the load balancer for the app should be enabled or disabled. Defaults - `false` | true |
| tlb.port | Expose the app through this port to the outside world. This is a mandatory label. If this property is not present even though `tlb.enabled` is set to `true` we'll not expose the app. | 11000 |
| tlb.portIndex | Ability to choose which port to be picked for load balancing. If you've configured 2 ports for your app, and you want to expose the app via the 2nd port via `tlb.port` to the outside world set this value to `1`. This is `0` based index value. Default - 0 | 1 |

## Contribute
If you've any feature requests or issues, please open a Github issue. We accept PRs. Fork away!

## License
http://www.apache.org/licenses/LICENSE-2.0