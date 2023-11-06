#
# 1. Build Container
#
FROM golang:1.18 as build

ENV GO111MODULE=on \
    GOOS=linux \
    GOPROXY=https://proxy.golang.org \
    GOARCH=amd64


# First add modules list to better utilize caching
COPY . /app
WORKDIR /app
# Download dependencies
RUN go mod download

# RUN go build

# Build components.
# Put built binaries and runtime resources in /app dir ready to be copied over or used.
# RUN go install -installsuffix cgo -ldflags="-w -s" && \
#     mkdir -p /app && \
#     cp -r $GOPATH/bin/go-echo-api /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /app


# 2. Runtime Container

FROM alpine

LABEL maintainer="Hector Nguyen <hoangnh2912@gmail.com>"

ENV TZ=Asia/Tokyo \
    PATH="/app:${PATH}"

RUN apk add --update --no-cache \
    tzdata \
    ca-certificates \
    bash \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone

# See http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /app

COPY --from=build /app/ /app/
EXPOSE 8080
CMD [ "./go-echo-api" ]