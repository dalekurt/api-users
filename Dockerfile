FROM golang:1.15.6-alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    API_HOME=/app

# 10001
# ARG GROUP_ID 
# 10000
# ARG USER_ID 

RUN apk --no-cache add ca-certificates git bind-tools

WORKDIR $API_HOME
# RUN go get -d -v golang.org/x/net/html

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o users .

FROM alpine:latest  

# Set necessary environmet variables needed for our image
ENV API_USER=nonroot \
    API_PORT=3000 \
    API_HOME=/app \
    GROUP_ID=10001 \
    USER_ID=10000

RUN apk --no-cache add ca-certificates
WORKDIR $API_HOME
COPY --from=0 $API_HOME/users .

# Non-root user for security purposes.
#
# UIDs below 10,000 are a security risk, as a container breakout could result
# in the container being ran as a more privileged user on the host kernel with
# the same UID.
#
# Static GID/UID is also useful for chown'ing files outside the container where
# such a user does not exist.
RUN addgroup -g $GROUP_ID -S nonroot && adduser -u $USER_ID -S -G nonroot -h /home/nonroot nonroot

# Install packages here with `apk add --no-cache`, copy your binary
# into /sbin/, etc.

# Tini allows us to avoid several Docker edge cases, see https://github.com/krallin/tini.
# RUN apk add --no-cache tini
# ENTRYPOINT ["/sbin/tini", "--", "/app/users"]
# Replace "myapp" above with your binary

# bind-tools is needed for DNS resolution to work in *some* Docker networks, but not all.
# This applies to nslookup, Go binaries, etc. If you want your Docker image to work even
# in more obscure Docker environments, use this.
RUN apk add --no-cache bind-tools

# Use the non-root user to run our application
# USER nonroot
USER $API_USER

EXPOSE $API_PORT

# Default arguments for your app (remove if you have none):
CMD ["./users"]
