FROM openresty/openresty AS l7bench-server-docker

LABEL maintainer="Sudarshan Raghavan <camelinx@gmail.com>"

ENV OPENRESTY_NGINX_BASE=/usr/local/openresty/nginx
ENV OPENRESTY_NGINX_CONF=$OPENRESTY_NGINX_BASE/conf
ENV OPENRESTY_NGINX_LUALIB=$OPENRESTY_NGINX_BASE/lualib
ENV OPENRESTY_NGINX_CERT_DIR=/etc/nginx-certs

COPY start-server.sh /root/

RUN chmod +x /root/start-server.sh && \
    mkdir -p $OPENRESTY_NGINX_CONF/conf.d && \
    mv $OPENRESTY_NGINX_CONF/nginx.conf $OPENRESTY_NGINX_CONF/nginx.conf.orig && \
    mkdir -p $OPENRESTY_NGINX_LUALIB && \
    mkdir -p $OPENRESTY_NGINX_CERT_DIR

COPY nginx.conf $OPENRESTY_NGINX_CONF/
COPY docker/l7bench-server.conf $OPENRESTY_NGINX_CONF/conf.d/
COPY docker/server.crt docker/server.key $OPENRESTY_NGINX_CERT_DIR/
COPY lualib/l7bench.lua $OPENRESTY_NGINX_LUALIB/

ENTRYPOINT ["/root/start-server.sh"]
