# syntax=docker/dockerfile:1
FROM golang:1.22.0-alpine3.19 as base



# ---------- LABELS
LABEL \
  maintainer="aleroxac" \
  vendor="aleroxac" \
  org.label-schema.author-name="Augusto Cardoso dos Santos" \
  org.label-schema.author-username="aleroxac" \
  org.label-schema.author-title="SRE/Platform Engineer" \
  org.label-schema.author-email="acardoso.ti@gmail.com" \
  org.label-schema.license="Apache-2.0" \
  org.label-schema.schema-version="1.0.0-rc.1" \
  org.label-schema.vcs-ref="" \
  org.label-schema.vcs-url="https://github.com/aleroxac/goexpert-weather-api" \
  org.label-schema.build-date="2024-02-27T09:06:56Z" \
  org.label-schema.name="goexpert-weather-api" \
  org.label-schema.base="golang:1.22.0-alpine3.19" \
  org.label-schema.version="v1" \
  org.label-schema.description="API to get the current weather based on Brazilian CEP city locations" \
  org.label-schema.usage="not-applicable" \
  org.label-schema.url="https://hub.docker.com/r/aleroxac/goexpert-weather-api:v1" \
  org.label-schema.os-name="not-applicable" \
  org.label-schema.os-version="not-applicable" \
  org.label-schema.docker.cmd="docker run --name=goexpert-weather-api --rm -d -p 8080:8080 aleroxac/goexpert-weather-api:v1" \
  org.label-schema.docker.cmd.devel="not-applicable" \
  org.label-schema.docker.cmd.test="not-applicable" \
  org.label-schema.docker.cmd.debug="not-applicable" \
  org.label-schema.docker.cmd.help="not-applicable" \
  org.label-schema.docker.params="not-applicable"



# ---------- ENVS
ENV DOCKER_CONTENT_TRUST=1



# ---------- BUILD
FROM base as build
WORKDIR /app
COPY go.mod main.go /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-api



# ---------- MAIN
FROM scratch
WORKDIR /app
COPY --from=build /app/weather-api .
ENTRYPOINT [ "./weather-api" ]
