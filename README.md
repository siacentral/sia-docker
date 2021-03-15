# Sia - Docker

[![Docker Pulls](https://img.shields.io/docker/pulls/siacentral/sia?color=19cf86&style=for-the-badge)](https://hub.docker.com/r/siacentral/sia)

An unofficial docker image for Sia. Automatically builds Sia using the source code from the official repository: https://gitlab.com/NebulousLabs/Sia

### Breaking change with Sia v1.5.6
With the Sia v1.5.6 update two potentially breaking changes will be made to this container: 
+ The `SIA_MODULES` environment variable will be removed, this was left in
largely for compatibility with `mtlynch/sia`. Instead you should pass `-M gct` directly at the end of `docker run` or as `command: -M gct` in docker-compose. 
+ siac and siad have been moved to `/usr/local/bin` for easier usage

# Release Tags

+ latest - the latest stable Sia release
+ beta - the latest release candidate for the next version of Sia
+ versions - builds of exact Sia releases such as: `1.5.4` or `1.5.5`
+ unstable - an unstable build of Sia's current master branch.

**Get latest official release:**
```
docker pull siacentral/sia:latest
```

**Get latest release candidate:**
```
docker pull siacentral/sia:beta
```

**Get Sia v1.5.4**
```
docker pull siacentral/sia:1.5.4
```

**Get unstable dev branch**
```
docker pull siacentral/sia:unstable
```

# Usage

It is important to never publish port `9980` to anything but 
`127.0.0.1:9980` doing so could give anyone full access to the Sia API and your
wallet.

Containers should never share volumes. If multiple sia containers are 
needed one unique volume should be created per container.

## Basic Container
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

### Command Line Flags

Additional siad command line flags can be passed in by appending them to docker
run.

#### Change API port from 9980 to 8880
```
docker run \
	--detach
	--restart unless-stopped \
	--publish 127.0.0.1:8880:8880 \
	--publish 9981:9981 \
	--publish 9982:9982 \
	--publish 9983:9983 \
	siacentral/sia --api-addr ":8880"
 ```


#### Change SiaMux port from 9983 to 8883
```
docker run \
	--detach
	--restart unless-stopped \
	--publish 127.0.0.1:9980:9980 \
	--publish 9981:9981 \
	--publish 9982:9982 \
	--publish 8883:8883 \
	siacentral/sia --siamux-addr ":8883"
 ```

#### Only run the minimum required modules
 ```
docker run \
	--detach
	--restart unless-stopped \
	--publish 127.0.0.1:9980:9980 \
	--publish 9981:9981 \
	--publish 9982:9982 \
	siacentral/sia -M gct
 ```

## Docker Compose

```yml
services:
  sia:
    container_name: sia
    image: siacentral/sia:latest
    ports:
      - 127.0.0.1:9980:9980
      - 9981:9981
      - 9982:9982
      - 9983:9983
      - 9984:9984
    volumes:
      - sia-data:/sia-data
    restart: unless-stopped

volumes:
  sia-data:
```

#### Change API port from 9980 to 8880
```yml
services:
  sia:
    container_name: sia
    command: --api-addr :8880
    image: siacentral/sia:latest
    ports:
      - 127.0.0.1:8880:8880
      - 9981:9981
      - 9982:9982
      - 9983:9983
      - 9984:9984
    volumes:
      - sia-data:/sia-data
    restart: unless-stopped

volumes:
  sia-data:
```


#### Change SiaMux port from 9983 to 8883
```yml
services:
  sia:
    container_name: sia
    command: --siamux-addr :8883
    image: siacentral/sia:latest
    ports:
      - 127.0.0.1:9980:9980
      - 9981:9981
      - 9982:9982
      - 8883:8883
      - 9984:9984
    volumes:
      - sia-data:/sia-data
    restart: unless-stopped

volumes:
  sia-data:
```

#### Only run the minimum required modules
```yml
services:
  sia:
    container_name: sia
    command: -M gct
    image: siacentral/sia:latest
    ports:
      - 127.0.0.1:9980:9980
      - 9981:9981
      - 9982:9982
      - 9983:9983
      - 9984:9984
    volumes:
      - sia-data:/sia-data
    restart: unless-stopped

volumes:
  sia-data:
```

## Sia API Password

When you create or update the Sia container a random API password will be
generated. You may need to copy the new API password when connecting outside of
the container. To force the same API password to be used you can add
`-e SIA_API_PASSWORD=yourpasswordhere` to the `docker run` command. This will
ensure that the API password stays the same between updates and restarts.

## Using Specific Modules

You can pass in different combinations of Sia modules to run by modifying the 
command used to create the container. For example: `-M gct` tells Sia to only
run the gateway, consensus, and transactionpool modules. `-M gctwh` is the minimum
required modules to run a Sia host. `-m gctwr` is the minimum required modules to
run a Sia renter.

## Hosts

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
