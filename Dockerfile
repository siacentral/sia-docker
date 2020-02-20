# build sia
FROM golang:1.13-alpine AS buildgo

ARG SIA_VERSION=master

WORKDIR /app

RUN echo "Install Build Tools" && apk update && apk upgrade && apk add --no-cache gcc musl-dev openssl git make

RUN echo "Clone Sia Repo" && git clone -b $SIA_VERSION https://gitlab.com/NebulousLabs/Sia.git /app

# docker makes GIT_DIRTY from the make file break even with a fresh repo
# updates git's index and makes it work properly again
RUN git diff --quiet; exit 0
RUN echo "Build Sia" && make release

# run sia
FROM alpine:latest

ENV SIA_MODULES gctwhr

EXPOSE 9980 9981 9982

COPY --from=buildgo /go/bin/sia* ./

ENTRYPOINT ./siad \
	--disable-api-security \
	-d /sia-data \
	--modules $SIA_MODULES \
	--api-addr ":9980"