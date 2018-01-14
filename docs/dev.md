# Development environment setup

Development environment relies on `docker` and `docker-compose` to have all services isolated, configured and to provide all required external dependencies (S3, load balancers, …).

## Location based setup

The development environment is setup to provide interface for two separate environments with their own domains:

* local clinic setup (`iryo.local`)
* cloud based setup (`iryo.cloud`)

Local development environment should be setup to forward all local traffic for both `iryo.local` and `iryo.cloud` to `127.0.0.1` or any other suitable IP address.

## Load balancing

Requests made on either of the two domains are first intercepted by `traefik`. It contains custom service definitions stored in the `/traefik.toml` configuration file that define how specific domains and paths should be forwarded to specific services.

Even though traefik supports `docker backend` we are not able to achieve fully encrypted communication between the traefik and services as certificates require communication with hostnames while traefik is only able to route traffic with specific IP addresses that it collects from docker. To work around this issue all services have to be defined manually with the `file backend`.

## TLS certficates

TLS certificates are generated for each server, utility or client. More details can be found [here](tls.md).