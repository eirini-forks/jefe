FROM golang:1.15 as builder
WORKDIR /build
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux go build -mod vendor -trimpath -a -installsuffix cgo -o jefe ./cmd/web/

FROM cloudfoundry/run:tiny
USER 2000:2000
COPY --from=builder /build/jefe .
COPY ./ui ./ui
ENTRYPOINT [ "./jefe" ]
