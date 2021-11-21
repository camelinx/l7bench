FROM openresty/openresty

LABEL maintainer="Sudarshan Raghavan <camelinx@gmail.com>"

ENV OPENRESTY_NGINX_BASE=/usr/local/openresty/nginx
ENV OPENRESTY_NGINX_CONF=$OPENRESTY_NGINX_BASE/conf
ENV OPENRESTY_NGINX_LUALIB=$OPENRESTY_NGINX_BASE/lualib

COPY start-server.sh /root/

RUN chmod +x /root/start-server.sh && \
    mkdir -p $OPENRESTY_NGINX_CONF/conf.d && \
    mv $OPENRESTY_NGINX_CONF/nginx.conf $OPENRESTY_NGINX_CONF/nginx.conf.orig && \
    mkdir -p $OPENRESTY_NGINX_LUALIB

COPY nginx.conf $OPENRESTY_NGINX_CONF/
COPY lualib/l7bench.lua $OPENRESTY_NGINX_LUALIB/

ENTRYPOINT ["/root/start-server.sh"]
