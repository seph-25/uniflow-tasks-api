# ---------- Fase de build (compila tu app) ----------
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Binario estático linux/amd64 -> ideal para contenedores
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# ---------- Fase de run (imagen final, mínima y segura) ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /src/app /app/app
EXPOSE 8080
# Usuario no root por seguridad
USER nonroot:nonroot
# Tu código usa PORT si existe; por defecto 8080
ENV PORT=8080
ENTRYPOINT ["/app/app"]
