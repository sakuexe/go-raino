FROM golang:alpine

ENV PROJECT raino
WORKDIR /go/src/${PROJECT}

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

COPY . .

# Download libwebp (for webp conversion)
RUN apk add build-base
RUN apk add --no-cache libwebp-dev

# Build the project
RUN CGO_ENABLED=1 GOOS=linux \
go build -o /usr/local/bin/${PROJECT} ./cmd

CMD ${PROJECT}
