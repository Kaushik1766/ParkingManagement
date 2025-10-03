FROM golang:1.25.1-alpine3.22

WORKDIR /app

ENV DATABASE_URL="postgresql://kaushik:123@postgres:5432/ParkingManagement"


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /web ./cmd/web/main.go

EXPOSE 3000

CMD ["/web"]
