This is an in-memory key-value store HTTP API service, with the following endpoints:

- `/get/{key}` : GET method. Returns the value of a previously set key.

- `/set` : POST method. Sets key/value pair(s) in the key-value store. Body can accept multiple key-value pairs in a single request, for example - `{"abc-1":1,"abc-2":2,"xyz-1":"three","xyz-2":4}`.

- `/search` : GET method. Searches for keys using prefix or suffix filters.

Assume you have the following keys in the store: abc-1, abc-2, xyz-1, xyz-2
`/search?prefix=abc` would return `abc-1` and `abc-2`.
`/search?suffix=-1` would return `abc-1` and `xyz-1`.

- `/` : GET method. Returns all the key-value pairs in the in-memory store. Also useful for readiness probes and load testing. 

## How to run (as binary)

- Fork and clone this repository.
- Run `go build -o kvstore`.
- Run the `./kvstore` binary (the service will run on port 8080). 

(Alternatively, you can simply run `go build` and run the resulting `./kv` binary, or give another name to it.)

## How to build and use Docker image

- Fork and clone this repository.
- Build the Docker image with the command
```
docker build -t <image-name>:<tag-name> .
```

For example, you can run `docker build -t kv-store:latest .`

- After building the image, run it with the command
```
docker run -d -p 8080:8080 <image-name>
```
(The service will run on port 8080 in this case.)

Alternatively, you can simply pull the docker image (the `latest` image tag is recommended) from the [Dockerhub repository](https://hub.docker.com/repository/docker/importhuman/kv-store).

## Running the service on a k3d cluster

- Fork and clone this repository (or simply copy the `deploy.yaml` file and save it).
- Create a new k3d cluster with `k3d cluster create kvstore -p "8080:80@loadbalancer"`.
(k3s [deploys traefik](https://k3d.io/v5.0.1/usage/exposing_services/) as the default ingress controller. Coupled with the port mapping, this does not require any further configuration for running on k3d. The deployment and configuration have not been tested for other Kubernetes distributions.)
- Run `kubectl apply -f deploy.yaml`.

The service will run on port 8080 in this case.

### Known issues

- While running on Kubernetes, the deployment creates 3 replicas of the key-value store service. These replicas do not have a shared memory, and so have different key-value pair stores, instead of all replicas utilizing a common store.