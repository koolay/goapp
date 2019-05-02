ARG GO_VERSION=1.11

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder
LABEL maintainer="koolay <koolay@163.com>"

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN \
    sed -i 's/dl-cdn.alpinelinux.org/mirror.tuna.tsinghua.edu.cn/g' /etc/apk/repositories; \
    apk add --no-cache ca-certificates curl build-base

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /go/src/github.com/koolay/goapp

# Import the code from the context.
COPY ./ ./

# Build the executable to `/app`. Mark the build as statically linked.
RUN OUTPUT=/app/goapp make; \
    cp gateway.json /app/

# Final stage: the running container.
FROM alpine:3.8 AS final

RUN apk add --no-cache curl ca-certificates

ENV APP_HOME=/home/app/webapp

# grab gosu for easy step-down from root
ENV GOSU_VERSION 1.11

ADD https://github.com/tianon/gosu/releases/download/${GOSU_VERSION}/gosu-amd64 /usr/local/sbin/gosu

ADD docker-entrypoint.sh /

# Import the compiled executable from the first stage.
COPY --from=builder /app/goapp $APP_HOME/goapp
COPY --from=builder /app/gateway.json $APP_HOME/gateway.json

RUN \
    sed -i 's/dl-cdn.alpinelinux.org/mirror.tuna.tsinghua.edu.cn/g' /etc/apk/repositories; \
    apk add --no-cache wget curl bash ca-certificates; \
	chmod +x /usr/local/sbin/gosu; \
	gosu nobody true

# force encoding
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

ARG UID=1000
ARG GID=1000

RUN \
    chmod 0755 /usr/local/sbin/gosu ;\
    chown root:root /usr/local/sbin/gosu ; \
    addgroup -g ${GID} app && \
    adduser -D -u ${UID} -s /bin/bash -G app app && \
    touch ${APP_HOME}/app.toml; \
    chown -R app:app ${APP_HOME}; \
    ln -s ${APP_HOME}/goapp /usr/local/bin/goapp; \
    chmod +x /docker-entrypoint.sh; \
    chmod +x /usr/local/bin/goapp; \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime; \
    echo "Asia/Shanghai" > /etc/timezone 

WORKDIR ${APP_HOME}
# Perform any further action as an unprivileged user.
EXPOSE 6666
# Run the compiled binary.
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["goapp", "serve"]
