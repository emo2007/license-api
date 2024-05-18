FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/app cmd/main.go

EXPOSE 8080

CMD ["/app/app", "-log-level=debug","-log-local=false","-log-add-source=true","-rest-address=0.0.0.0:8080","-db-host=blockd-db:5432","-db-database=blockd","-db-user=blockd","-db-secret=blockd","-db-enable-tls=false"]
