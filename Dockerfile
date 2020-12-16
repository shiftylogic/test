FROM golang:1.15-alpine3.12 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o site ./cmd


FROM scratch
COPY --from=builder /build/site /app/
WORKDIR /app
CMD ["./site"]

