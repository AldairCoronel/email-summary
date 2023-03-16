FROM ubuntu:latest

# Set noninteractive mode
ENV DEBIAN_FRONTEND=noninteractive

# Install Go 1.17 and PostgreSQL
RUN apt-get update && \
    apt-get install -y golang-1.17-go postgresql ca-certificates

# Set environment variables
ENV GOROOT=/usr/lib/go-1.17
ENV GOPATH=/app/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Install PostgreSQL
RUN apt-get install -y postgresql

# Create workdir for project
RUN mkdir /app

# Set working directory to project directory
WORKDIR /app

# Copy project files to container
COPY . .

# Download dependencies
RUN go mod download

# Give path to project
ENV PATH_TO_PROJECT=/app

# Expose port for PostgreSQL
EXPOSE 5432

# Start the PostgreSQL server
CMD service postgresql start && bash
