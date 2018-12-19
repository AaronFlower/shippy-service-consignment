# We use the official golang image, which contains all the
# current build tools and libraries. 
# Notice `as builder`, this gives the container a name that
# we can reference later on.

FROM golang:1.11 as builder

# Set out workdir to our current service in the gopath
WORKDIR /go/src/github.com/aaronflower/shippy-service-consignment

# Copy the current code into our workdir
COPY . .

# Here we're pulling in godep, which is a dependency manager tool,
# we're going to use dep instead of go get, to get around a few
# quirks in how go get works with sub-packages.
# RUN go get -u -v github.com/golang/dep/cmd/dep
# RUN dep init -v && dep ensure -v

# Build the binary, with a few flags which will allow us to run this binary in Alpine.
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

# Here we're using a second FROM statement, which is strange,
# but this tells Docker to start a new build process with this image.
FROM alpine:latest

# Security related package, good to have.
RUN apk --no-cache add ca-certificates

# Same as before, create a directory for our app.
RUN mkdir /app
WORKDIR /app

# Here, instead of copying the binary from our host machine, we pull the binary from
# the container name `builder`, within this build context. This reaches into our previous
# image, finds the binary we built, and pulls it into this container. Amazing.
COPY ./service.consignment .

# Run the binary as per usual! This time with a binary build in a separate container,
# with all of the correct dependencies and runtime libraries.
CMD ["./service.consignment"]

