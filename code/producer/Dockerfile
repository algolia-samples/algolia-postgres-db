FROM golang:1.18.3-alpine as build
RUN apk update
RUN apk add -U --no-cache ca-certificates && update-ca-certificates
WORKDIR /appbuild
COPY . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o audit_producer .

FROM scratch as run
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /appbuild/audit_producer /app/audit_producer

ENTRYPOINT [ "/app/audit_producer" ]

