FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN /bin/sh -c 'CGO_ENABLED=0 GOOS=linux go build -o /bin/handler .'


FROM debian:11-slim

WORKDIR /app

RUN apt-get update \
    && apt-get install -y curl gnupg2 lsb-release libicu-dev \
    && apt-get clean

RUN curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.gpg
RUN mv microsoft.gpg /etc/apt/trusted.gpg.d/microsoft.gpg

RUN sh -c 'echo "deb [arch=amd64] https://packages.microsoft.com/debian/$(lsb_release -rs | cut -d'.' -f 1)/prod $(lsb_release -cs) main" > /etc/apt/sources.list.d/dotnetdev.list'

RUN apt-get update \
    && apt-get install -y azure-functions-core-tools-4 \
    && apt-get clean

COPY host.json ./

COPY listings/function.json ./listings/
COPY listingById/function.json ./listingById/

COPY --from=builder /bin/handler ./bin/

EXPOSE 80

ENTRYPOINT [ "func", "start", "--custom", "--port", "80", "--verbose" ]
