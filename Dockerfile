# -------------------------
# Build Stage
# -------------------------
FROM golang:1.25.3-bookworm AS builder

WORKDIR /ntp

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags='-s -w' -o ntp .

# -------------------------
# Final Stage (scratch)
# -------------------------
FROM scratch

WORKDIR /ntp

COPY --from=builder /ntp/ntp .

CMD ["/ntp/ntp"]