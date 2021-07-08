# The Go version that I have made the project on. This is the base image
FROM golang:1.16.5

# Set the cwd in the container
WORKDIR /app

# Get the port where Go will run. We will get this from docker-compose
ARG API_PORT

# Copy mod and sum file. These files will be used in the container to 
# download the dependenciesßß
COPY ./go.mod ./go.sum ./

# Download Go dependencies
RUN go mod download

# Copy host files in the container
COPY . .

# Expose go port
EXPOSE ${API_PORT}