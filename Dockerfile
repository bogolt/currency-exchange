# # Start from a Debian image with the latest version of Go installed
# # and a workspace (GOPATH) configured at /go.
# FROM golang

# # Copy the local package files to the container's workspace.
# ADD . /go/src/github.com/bogolt/currency-exchange

# # Build the outyet command inside the container.
# # (You may fetch or manage dependencies here,
# # either manually or with a tool like "godep".)
# RUN go install github.com/bogolt/currency-exchange

# # Run the outyet command by default when the container starts.
# ENTRYPOINT /go/bin/currency-exchange

# FROM golang:onbuild
# # Document that the service listens on port 8080.
# EXPOSE 8080


FROM golang:1.15-alpine

# Create a directory for the app
RUN mkdir /app
 
# Copy all files from the current directory to the app directory
COPY . /app
 
# Set working directory
WORKDIR /app
 
# Run command as described:
# go build will build an executable file named server in the current directory
RUN go build -o currency-exchnage . 
 
# Run the server executable
CMD [ "/app/currency-exchnage" ]