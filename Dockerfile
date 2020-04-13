# build sia
FROM golang:1.13-alpine AS buildgo

ARG SIA_VERSION=master

RUN echo "Install Build Tools" && apk update && apk upgrade && apk add --no-cache gcc musl-dev openssl git make

# prevents cache on git clone if the ref has changed
ADD https://gitlab.com/api/v4/projects/7508674/repository/commits/${SIA_VERSION} version.json

WORKDIR /app

RUN echo "Clone Sia Repo" && git clone -b $SIA_VERSION https://gitlab.com/NebulousLabs/Sia.git /app

# docker makes GIT_DIRTY from the make file break even with a fresh repo
# updates git's index and makes it work properly again
RUN git diff --quiet; exit 0

RUN echo "Build Sia" && mkdir /app/releases && go build -a -tags 'netgo' -trimpath \
	-ldflags="-s -w -X 'gitlab.com/NebulousLabs/Sia/build.GitRevision=`git rev-parse --short HEAD`' -X 'gitlab.com/NebulousLabs/Sia/build.BuildTime=`date`'" \
	-o /app/releases ./cmd/siad ./cmd/siac

# run sia
FROM alpine:latest

ENV SIA_MODULES gctwhr

EXPOSE 9980 9981 9982

COPY --from=buildgo /app/releases ./

ENTRYPOINT ./siad \
	--disable-api-security \
	-d /sia-data \
	--modules $SIA_MODULES \
	--api-addr ":9980"