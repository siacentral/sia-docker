# build sia
FROM golang:1.13-alpine AS buildgo

ARG SIA_VERSION=master

WORKDIR /app

RUN echo "Install Build Tools" && apk update && apk upgrade && apk add --no-cache gcc musl-dev openssl git

RUN echo "Clone Sia Repo" && git clone -b $SIA_VERSION https://gitlab.com/NebulousLabs/Sia.git /app

RUN echo "Build Sia" && mkdir /app/releases && go build -a -tags 'netgo' -trimpath \
	-ldflags="-s -w -X 'gitlab.com/NebulousLabs/Sia/build.GitRevision=`git rev-parse --short HEAD`' -X 'gitlab.com/NebulousLabs/Sia/build.BuildTime=`date`' -X 'gitlab.com/NebulousLabs/Sia/build.ReleaseTag=${SIA_VERSION}'" \
	-o /app/releases ./cmd/siad ./cmd/siac

# run sia
FROM alpine:latest

ENV SIA_MODULES gctwhr

EXPOSE 9980 9981 9982

COPY --from=buildgo /app/releases/* ./

RUN echo "Install Socat" && apk update && apk upgrade && apk add --no-cache socat

ENTRYPOINT socat tcp-listen:9980,reuseaddr,fork tcp:localhost:8000 & \
	./siad -d /sia-data --modules gctwhr --api-addr "localhost:8000"