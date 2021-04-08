# build sia
FROM golang:1.15-alpine AS buildgo

ARG SIA_VERSION=master
ARG RC=master

RUN echo "Install Build Tools" && apk update && apk upgrade && apk add --no-cache gcc musl-dev openssl git make

# prevents cache on git clone if the ref has changed
ADD https://gitlab.com/api/v4/projects/7508674/repository/commits/${SIA_VERSION} version.json

WORKDIR /app

RUN echo "Clone Sia Repo" && git clone https://gitlab.com/NebulousLabs/Sia.git /app && git fetch && git checkout $SIA_VERSION

RUN echo "Build Sia" && mkdir /app/releases && go build -a -tags 'netgo' -trimpath \
	-ldflags="-s -w -X 'gitlab.com/NebulousLabs/Sia/build.GitRevision=`git rev-parse --short HEAD`' -X 'gitlab.com/NebulousLabs/Sia/build.BuildTime=`git show -s --format=%ci HEAD`' -X 'gitlab.com/NebulousLabs/Sia/build.ReleaseTag=${RC}'" \
	-o /app/releases ./cmd/siad ./cmd/siac

# run sia
FROM alpine:latest

COPY --from=buildgo /app/releases /usr/local/bin

EXPOSE 9981 9982 9983 9984

ENV SIA_WALLET_PASSWORD=
ENV SIA_API_PASSWORD=

VOLUME [ "/sia-data" ]

ENTRYPOINT [ "siad", "--disable-api-security", "-d", "/sia-data", "--api-addr", ":9980" ]
