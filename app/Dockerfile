# syntax = docker/dockerfile:experimental
# Build Container
FROM golang:1.17.2 as builder

ENV GO111MODULE on
ENV GOPRIVATE "bitbucket.org/latonaio"

WORKDIR /go/src/bitbucket.org/latonaio

COPY ./app/go.mod .

RUN git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"

RUN mkdir /root/.ssh/ && touch /root/.ssh/known_hosts && ssh-keyscan -t rsa bitbucket.org >> /root/.ssh/known_hosts

RUN --mount=type=secret,id=ssh,target=/root/.ssh/id_rsa go mod download

COPY ./app .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o control-mitsubishi-plc-w-kube .


# Runtime Container
FROM alpine:3.12

RUN apk add --no-cache libc6-compat


ENV SERVICE=control-mitsubishi-plc-w-kube \
    POSITION=Runtime \
    AION_HOME="/var/lib/aion" \
    APP_DIR="${AION_HOME}/${POSITION}/${SERVICE}"

WORKDIR ${APP_DIR}

COPY --from=builder /go/src/bitbucket.org/latonaio/control-mitsubishi-plc-w-kube .

CMD ["./control-mitsubishi-plc-w-kube"]
