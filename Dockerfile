# -------------------------
# Build Stage
# -------------------------
FROM golang:1.25.3-bookworm AS builder

WORKDIR /ntp

# Only copy go.mod/go.sum for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build fully static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags='-s -w' -o ntp .

# -------------------------
# Final Stage (scratch)
# -------------------------
FROM scratch

WORKDIR /ntp

# Copy the statically compiled binary
COPY --from=builder /ntp/ntp .

# Run the binary
CMD ["/ntp/ntp"]