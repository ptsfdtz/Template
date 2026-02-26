# builder for react
FROM node:22-slim AS react-builder
WORKDIR /client
RUN corepack enable
COPY client/package.json client/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY client .
RUN pnpm build

# builder for go
FROM golang:1.26-alpine AS go-builder
WORKDIR /server
COPY server .
COPY --from=react-builder /client/dist/ build/
RUN go mod download
RUN mkdir -p bin
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bin/app ./main.go
RUN apk update && apk add --no-cache ca-certificates upx
RUN upx --best bin/app

# production image
FROM scratch
EXPOSE 3001
COPY --from=go-builder /server/bin/app /app
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app"]
