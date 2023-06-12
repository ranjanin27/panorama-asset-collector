# ------------------------------------------------------------------------------
# (C) Copyright 2022 Hewlett Packard Enterprise Development LP
# ------------------------------------------------------------------------------


# ------------------------------------------------------------------------------
# Shared base image for other containers
# ------------------------------------------------------------------------------
FROM cds-harbor.rtplab.nimblestorage.com/docker_proxy/library/golang:1.19-alpine3.16 as base
#FROM golang:1.17-alpine3.15 as base

ARG USER=default
#ENV http_proxy 'http://web-proxy.in.hpecorp.net:8080'
#ENV https_proxy 'http://web-proxy.in.hpecorp.net:8080'
# install build tools
RUN set -eux; \
    apk add -U --no-cache \
        curl \
        git  \
        make \
        bash \
	build-base \
    ;

# add new user
RUN addgroup -g 1000 ${USER} \
        && adduser -h /build -D -u 1000 -G ${USER} ${USER} \
    ;

RUN echo  ${USER}
USER ${USER}
WORKDIR /build

# ------------------------------------------------------------------------------
# Build image
# ------------------------------------------------------------------------------

FROM base as build

ARG USER=default
ARG TOKEN=""
COPY --chown=${USER}:${USER} . ./

ARG VERSION=0.0.0
# RUN git remote set-url origin 'https://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@github.hpe.com/pruthvi-raju/panorama-common-hauler.git'
# COPY .netrc ${HOME}/
# RUN  echo ${HOME}
# RUN set -eux; cat ${HOME}/.netrc; \
#  make build VERSION=${VERSION};
RUN set -eux; \
   if [[ ! -z "$TOKEN" ]]; then \
       git config --global url."https://x-oauth-basic:${TOKEN}@github.hpe.com/".insteadof https://github.hpe.com/; \
   fi; \
   go env -w GOPRIVATE="github.hpe.com"; \
   make build VERSION=${VERSION};

# ------------------------------------------------------------------------------
# Runtime image
# ------------------------------------------------------------------------------

FROM cds-harbor.rtplab.nimblestorage.com/docker_proxy/library/alpine:3.16 as runtime
#FROM golang:1.17-alpine3.15 as runtime

ARG USER=default
ENV HOME /home/${USER}

# add new user
RUN addgroup -g 1000 ${USER} \
        && adduser -D -u 1000 -G ${USER} ${USER} \
    ;

USER ${USER}
WORKDIR /usr/local/bin
COPY --from=build /build/dist/panorama-common-hauler ./panorama-common-hauler

ENTRYPOINT ["/usr/local/bin/panorama-common-hauler"]
