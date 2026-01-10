FROM golang:alpine AS builder

WORKDIR /english-speaking-club-bot

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/english-speaking-club-bot ./cmd/app

FROM alpine

COPY --from=builder /bin/english-speaking-club-bot /bin/english-speaking-club-bot

EXPOSE 80

ENTRYPOINT ["/bin/english-speaking-club-bot"]
