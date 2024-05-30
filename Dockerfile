# syntax=docker/dockerfile:1.4

FROM ubuntu:jammy AS build-dev
WORKDIR /go/src/app
COPY --link go.mod go.sum ./
RUN apt update && apt install -y wget gcc golang-1.21
RUN ln -s /usr/lib/go-1.21/bin/go /usr/bin/go
RUN go version && go mod download
COPY --link . .
RUN wget https://github.com/mattn/go-tflite/releases/download/v1.0.5/go-tflite-buildkit-20240529.tar.gz && \
    tar -C /usr/local -xvf go-tflite-buildkit-20240529.tar.gz && \
    rm go-tflite-buildkit-20240529.tar.gz
RUN CGO_ENABLED=1 go build -buildvcs=false -trimpath -ldflags '-w -s' -o /usr/bin/mnist
FROM ubuntu:jammy
COPY --link --from=build-dev /usr/bin/mnist /
COPY --link --from=build-dev /usr/local/lib/libtensorflowlite_c.so /usr/local/lib/libtensorflowlite_c.so
RUN ldconfig
COPY --link mnist_model.tflite /mnist_model.tflite
ENTRYPOINT ["/mnist"]
