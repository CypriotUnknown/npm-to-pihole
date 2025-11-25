# npm-to-pihole

> Sync Nginx Proxy Manager proxy hosts to Pi-hole CNAME records.

npm-to-pihole watches Nginx proxy host configuration files, extracts `server_name` entries, and keeps Pi-hole's CNAME records in sync by adding missing records and removing stale ones..

## Features

- Automatically adds CNAME records to Pi-hole for `server_name` entries found in Nginx Proxy Manager proxy host files.
- Removes CNAME records from Pi-hole when the corresponding proxy host is removed.
- Designed to run as a small Docker service and integrate with Docker Compose.

## How it works

> Note: The service communicates with Pi-hole using HTTPS and sets `InsecureSkipVerify` to `true` in the HTTP client (accepts self-signed certs). Run it only on a trusted internal network.

## Quick start with Docker Compose

Here is an example Docker Compose configuration to run the app:

Set the `PIHOLE_PASSWORD` environment variable to your Pi-hole web admin password. Example snippet from `compose.sample.yml`:

```yaml
---
services:
  npm_to_pihole:
    image: gitlab.com/cypriotunknown/npm-to-pihole
    restart: unless-stopped
    volumes:
      - nginx:/data/nginx:ro
    networks:
      - pihole
    environment:
      PIHOLE_PASSWORD: <PIHOLE APP PASSWORD>
      NGINX_PROXY_DIR: /data/nginx/nginx/proxy_host
```

## Environment variables

- `PIHOLE_PASSWORD` (required): The Pi-hole web/admin password used to authenticate with the Pi-hole API.
- `NGINX_PROXY_DIR` (optional): Path inside the container to the folder that contains Nginx Proxy Manager `proxy_host` configuration files. Default used by the program is `/data/nginx/nginx/proxy_host`.

## Running locally (for development)

You can build and run the Go binary locally if you want to test changes.

Requirements: Go 1.25+ (or the version used in the project)

```bash
# build
go build -o npm-to-pihole ./

# run (ensure PIHOLE_PASSWORD is set and mount your nginx proxy dir appropriately)
PIHOLE_PASSWORD=yourpassword NGINX_PROXY_DIR=/path/to/proxy_host ./npm-to-pihole
```

## Security considerations

- The HTTP client disables TLS verification for Pi-hole (accepts self-signed certificates). Only deploy this in trusted networks.
- Keep Pi-hole credentials safe. Using Docker secrets or an environment management solution is recommended for production.

## Troubleshooting

- If the service cannot authenticate, verify `PIHOLE_PASSWORD` and network connectivity to the Pi-hole host/service.
- If changes are not reflected, check that the `NGINX_PROXY_DIR` correctly points to the directory containing `proxy_host` config files and that files include `server_name` directives.
- Check logs with `docker compose logs -f npm_to_pihole` for detailed errors.

## Contributing

Contributions are welcome. Please open issues if you come accross bugs.