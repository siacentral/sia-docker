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
An unofficial docker image for Sia. Automatically builds Sia using the source code from the official repository: https://gitlab.com/siacentral/sia

# Release Tags

+ latest - the latest stable Sia release
+ versions - builds of exact Sia releases such as: `1.4.4`, `1.4.3` or `1.4.2.1`
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

### Basic Container
```
docker volume create sia-data
docker run \
  --detach \
  --restart unless-stopped \
  --mount type=volume,src=sia-data,target=/sia-data \
  --publish 127.0.0.1:9980:9980 \
  --publish 9981:9981 \
  --publish 9982:9982 \
  --name sia \
   siacentral/sia
```

It is important to never `--publish` port `9980` to anything but 
`127.0.0.1:9980` doing so could give anyone full access to the Sia API and
wallet.

`docker volume create sia-data` creates a new persistent volume called 
"sia-data" to store Sia's data and blockchain. This will allow for the 
blockchain to remain consistent between container restarts or updates.

Containers should never share volumes. If multiple sia containers are 
needed one unique volume should be created per container.

## Sia API Password

When you create or update the Sia container a random API password will be
generated. You may need to copy the new API password when connecting outside of
the container. To force the same API password to be used you can add
`-e SIA_API_PASSWORD=yourpasswordhere` to the `docker run` command. This will
ensure that the API password stays the same between updates and restarts.

### Using Specific Modules

By specifying the docker argument `-e` you can pass in different combinations of
Sia modules to run. For example: `-e SIA_MODULES="gct"` tells Sia to only run
the gateway, consensus, and transactionpool modules.

#### Consensus Only
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
  --name sia \
   siacentral/sia
```

#### Renter Only
```
docker volume create sia-data
docker run \
  --detach \
  --restart unless-stopped \
  -e SIA_MODULES="gctwr" \
  --mount type=volume,src=sia-data,target=/sia-data \
  --publish 127.0.0.1:9980:9980 \
  --publish 9981:9981 \
  --publish 9982:9982 \
  --name sia \
   siacentral/sia
```

#### Host Only
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
  --name sia \
   siacentral/sia
```

Hosting may require additional volumes passed into the container to map
local drives into the container. These can be added by specifying
docker's `-v` or `--mount` flag.

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
