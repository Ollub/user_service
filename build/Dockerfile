# Step 0 - collect dependencies
FROM golang:1.18 AS dependencies
WORKDIR /go/src/user_service
COPY go.mod .
COPY go.sum .
RUN go mod download

# Step 1 - build
FROM dependencies AS build
COPY . /go/src/user_service
WORKDIR /go/src/user_service
RUN make build

# Step 2 - application image
FROM debian:buster-slim
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/user_service/bin/user_service /app/
COPY --from=build /go/src/user_service/scripts/wait-for-it.sh /app/
RUN chmod +x /app/*
#EXPOSE 8080/tcp
CMD /app/user_service
