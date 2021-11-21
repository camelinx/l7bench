#!/bin/sh

OPENRESTY_BASE=/usr/local/openresty

if [ -z ${NGX_CONFIGMAP_NAME} ]; then
    echo ${NGX_CONFIGMAP_NAME} >${OPENRESTY_BASE}/nginx/conf/conf.d/sample.conf
fi

${OPENRESTY_BASE}/bin/openresty -g "daemon off;"
