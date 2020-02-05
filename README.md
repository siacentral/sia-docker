Based on https://github.com/mtlynch/docker-sia but allows for building any commit or
version of Sia from source instead of downloading the official releases.

# Usage

## Basic Full Node

```
docker volume create sia-data
docker run \
  --detach \
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
  --env gct \
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
  --env gctwr \
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
  --env gctwh \
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
docker build --build-arg SIA_VERSION=v1.4.2.1 -t siacentral/docker-sia .
```