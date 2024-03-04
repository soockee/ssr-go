# ssr-go


## Deploy:

docker run --rm -v ./certs:/certs -p 443:443 -p 80:80 ghcr.io/soockee/ssr-go:latest -isProd