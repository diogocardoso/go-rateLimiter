# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependência
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o rateLimiter .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copiar o binário compilado do stage anterior
COPY --from=builder /app/rateLimiter .
COPY --from=builder /app/.env .

# Expor a porta da aplicação
EXPOSE 8080

# Executar a aplicação
CMD ["/app/rateLimiter"]