# Deployer CLI Tool

## Prerequisites

1. Install Pre-requisites:
   - [Docker](https://docs.docker.com/get-docker/)
   - [Go](https://go.dev/doc/install) (version 1.20 or later)
   - [Helm](https://helm.sh/docs/intro/install/) (version 3.0 or later)
   - [kubectl](https://kubernetes.io/docs/tasks/tools/) (version 1.20 or later)
   - [k9s](https://github.com/derailed/k9s) (version 0.25 or later)
2. Install Dependencies
    ```bash
    go mod download
    ```

2. Build the CLI tool
    ```bash
    go build cmd/main/deployer.go
    ```
   
## Run the CLI tool

**Deploy the k3s cluster**

```bash
./deployer cluster
```

or using Docker:

```bash
docker run --rm \
  --network host \
  --privileged \
  -v /etc/hosts:/etc/hosts \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /tmp:/tmp \
  mikeengineering/deployer \
  cluster
```

**Deploy your application's helm chart**

```bash
./deployer deploy --helm oci://registry-1.docker.io/repository/chart -n namespace -f values.yaml
```

or using Docker:

```bash
docker run --rm \
  --network host \
  --privileged \
  -v /etc/hosts:/etc/hosts \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd)/values.yaml:/app/values.yaml \
  -v /tmp:/tmp \
  mikeengineering/deployer \
  deploy --helm oci://registry-1.docker.io/repository/chart -n namespace -f values.yaml
```

**Check the status of the application**

```bash
./deployer status
```

**Open the application**

http://localhost:30080/

## Cleanup

**Delete the k3s cluster**

```bash
./deployer destroy
```

Or using Docker:

```bash
docker run --rm \
  --network host \
  --privileged \
  -v /etc/hosts:/etc/hosts \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /tmp:/tmp \
  mikeengineering/deployer \
  destroy
```
