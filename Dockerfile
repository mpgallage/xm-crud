# use official golang image as base image
FROM golang:1.20.4-bullseye

# set working directory
WORKDIR /app

# copy the source from the current directory to the working directory inside the container
COPY . .

# download dependencies
RUN go get -d -v ./...

# run linter
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2
RUN golangci-lint run

# build the go app
RUN go build -o github.com/mpgallage/xmcrud .

# expose port 8080 to the outside world
EXPOSE 8080

# command to run the executable
CMD ["./github.com/mpgallage/xmcrud"]

