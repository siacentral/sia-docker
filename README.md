# Sia - Docker

[![Docker Pulls](https://img.shields.io/docker/pulls/siacentral/sia?color=19cf86&style=for-the-badge)](https://hub.docker.com/r/siacentral/sia)

An unofficial docker image for Sia. Automatically builds Sia using the source code from the official repository: https://gitlab.com/NebulousLabs/Sia

# Release Tags

+ latest - the latest stable Sia release
+ beta - the latest release candidate for the next version of Sia
+ versions - builds of exact Sia releases such as: `1.4.7` or `1.4.6`
+ unstable - an unstable build of Sia's current master branch.

**Get latest official release:**
```
docker pull siacentral/sia:latest
```

**Get latest release candidate:**
```
docker pull siacentral/sia:beta
```

**Get Sia v1.4.5**
```
docker pull siacentral/sia:1.4.5
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
  --publish 9983:9983 \
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

### Sia API Password

When you create or update the Sia container a random API password will be
generated. You may need to copy the new API password when connecting outside of
the container. To force the same API password to be used you can add
`-e SIA_API_PASSWORD=yourpasswordhere` to the `docker run` command. This will
ensure that the API password stays the same between updates and restarts.

### Command Line Flags

Additional siad command line flags can be passed in by appending `-c` to docker run.


#### Change SiaMux port from 9983 to 8883
```
docker run \
  --detach
  --restart unless-stopped \
 --publish 127.0.0.1:9980:9980 \
 --public 9981:9981 \
 --publish 9982:9982 \
 --publish 8883:8883 \
 siacentral/sia -c --siamux-addr ":8883"
 ```

#### Change Sia API user-agent from "Sia-Agent" to "Custom-Agent"
 ```
docker run \
  --detach
  --restart unless-stopped \
 --publish 127.0.0.1:9980:9980 \
 --public 9981:9981 \
 --publish 9982:9982 \
 siacentral/sia -c --agent "Custom-Agent"
 ```


### Using Specific Modules

By specifying the environment variable `SIA_MODULES` you can pass in different combinations of
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
  --publish 9983:9983 \
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
docker build --build-arg SIA_VERSION=v1.4.7 -t siacentral/sia:1.4.7 .
```
