FROM golang:alpine AS build

RUN apk add git

RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go build -o /tmp/server ./cmd/api/main.go

FROM alpine:edge

COPY --from=build /tmp/server /sbin/server

# CMD /sbin/server
# ENTRYPOINT [ "/sbin/server", "--port=$PORT", "--mongo-db-name=$MONGO_DB_NAME", "--mongo-db-uri=$MONGO_DB_URL" ]