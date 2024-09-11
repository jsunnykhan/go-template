FROM golang:1.22-alpine

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

COPY . .

RUN CGO_ENABLED=0 GOFLAGS=-mod=mod GOOS=linux  go build -a -o /app ./cmd/init


EXPOSE 8000

# Run the compiled binary.
ENTRYPOINT ["/app", "server"]
