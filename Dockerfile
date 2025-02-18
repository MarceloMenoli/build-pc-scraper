# Etapa de build
FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Etapa final - usando Alpine para imagem final
FROM alpine:3.16

# Instala o Chromium e dependências essenciais no Alpine
RUN apk update && apk add --no-cache \
  chromium \
  ca-certificates \
  ttf-freefont

WORKDIR /root/
COPY --from=builder /app/main .

# Define a variável de ambiente para o caminho do Chromium
ENV CHROME_PATH=/usr/bin/chromium-browser
# Em Alpine, o binário do chromium pode estar em /usr/bin/chromium ou /usr/bin/chromium-browser.
# Verifique com "which chromium" no container.

EXPOSE 8080
CMD ["./main"]
