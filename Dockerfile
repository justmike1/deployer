FROM golang:1.23 AS build

ENV CGO_ENABLED=0
ENV GOARCH=amd64

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /go/bin/deployer /app/cmd/main/deployer.go

RUN curl -sfL https://github.com/k3s-io/k3s/releases/latest/download/k3s -o /usr/local/bin/k3s && chmod +x /usr/local/bin/k3s
RUN curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
RUN curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

RUN apt-get update && apt-get install -y busybox-static

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/deployer /deployer
COPY --from=build /usr/local/bin/k3s /usr/local/bin/k3s
COPY --from=build /usr/local/bin/k3d /usr/local/bin/k3d
COPY --from=build /usr/local/bin/helm /usr/local/bin/helm
COPY --from=build /bin/busybox /bin/sh

ENTRYPOINT ["/deployer"]