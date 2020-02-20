An unofficial docker image for Sia. Automatically builds Sia using the source code from the official repository: https://gitlab.com/NebulousLabs/Sia

# Release Tags

+ latest - the latest stable Sia release
+ versions - builds of exact Sia releases such as: `1.4.3` or `1.4.2.1`
+ unstable - an unstable build of Sia's master branch. Updated every 15 minutes

**Get latest release:**
```
docker pull siacentral/sia:latest
```

**Get Sia v1.4.2.1**
```
docker pull siacentral/sia:1.4.2.1
```

**Get unstable dev branch**
```
docker pull siacentral/sia:unstable
```


# Usage

## Basic Full Node

```
docker volume create sia-data
docker run \
  --detach \
  --restart unless-stopped \
  --mount type=volume,src=sia-data,target=/sia-data \
  --publish 127.0.0.1:9980:9980 \
  --publish 9981:9981 \
  --publish 9982:9982 \
  --name sia-temp \
   siacentral/sia
```

## Consensus Only Node

```
docker volume create sia-data
docker run \
  --detach \
  --restart unless-stopped \
  -e SIA_MODULES="gct" \
  --mount type=volume,src=sia-data,target=/sia-data \
  --publish 127.0.0.1:9980:9980 \
  --publish 9981:9981 \
  --publish 9982:9982 \
  --name sia-temp \
   siacentral/sia
```

## Rent Only Node

```
docker volume create sia-data
docker run \
  --detach \
  -e SIA_MODULES="gctwr" \
  --mount type=volume,src=sia-data,target=/sia-data \
  --publish 127.0.0.1:9980:9980 \
  --publish 9981:9981 \
  --name sia-temp \
   siacentral/sia
```

## Host Only Node

```
docker volume create sia-data
docker run \
  --detach \
  --restart unless-stopped \
  -e SIA_MODULES="gctwh" \
  --mount type=volume,src=sia-data,target=/sia-data \
  --publish 127.0.0.1:9980:9980 \
  --publish 9981:9981 \
  --publish 9982:9982 \
  --name sia-temp \
   siacentral/sia
```

## Building

To build a specific commit or version of Sia specify the tag or branch of the 
repository using Docker's `--build-arg` flag. Any valid `git checkout` ref can
be used with the `SIA_VERSION` build arg.

```
docker build --build-arg SIA_VERSION=v1.4.2.1 -t siacentral/sia:1.4.2.1 .
```

## /build

A simple GoLang CLI that checks any tags matching the version ID regex
from NebulousLabs/Sia and compares them to matching tags from a Docker Hub repo. 
It automatically builds and pushes any missing versions of Sia.

Includes some logic for `latest` and `unstable` tags. This CLI is run on Sia 
Central's build server via `cron` every 15 minutes to keep it automatically
updated with Sia's latest releases.

**Build**

```
go install build/build.go
```

**Run**
```
build --docker-hub-repo siacentral/sia
```