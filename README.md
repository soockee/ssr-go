# ssr-go


## Temp deploy:

docker run --rm -it -v <>  -e DOMAIN_NAME=<> -e DEPLOYMENT_ENVIRONMENT=<> -e SSL_CACHE_DIR=<>  -p 433:443 -p 80:80 ghcr.io/soockee/ssr-go:1.0.11 