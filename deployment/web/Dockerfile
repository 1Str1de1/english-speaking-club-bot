FROM golang:alpine AS builder

WORKDIR /english-speaking-club-Bot

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/english-speaking-club-Bot ./cmd/app

FROM alpine

COPY --from=builder /bin/english-speaking-club-Bot /bin/english-speaking-club-Bot

EXPOSE 80

ENTRYPOINT ["/bin/english-speaking-club-Bot"]
