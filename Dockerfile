FROM golang:1.26-alpine AS api-build

WORKDIR /src/apps/api
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download
COPY apps/api/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/server

FROM node:22-alpine AS web-build

WORKDIR /src/apps/web
COPY apps/web/package.json apps/web/package-lock.json ./
RUN npm ci
WORKDIR /src
COPY apps/web/ ./apps/web/
COPY contracts/ ./contracts/
COPY data/ ./data/
ENV NEXT_PUBLIC_API_URL=/ NEXT_PUBLIC_USE_MOCKS=false
RUN npm --prefix apps/web run build

FROM node:22-alpine

WORKDIR /app
ENV NODE_ENV=production \
    HOSTNAME=0.0.0.0 \
    PORT=3000 \
    API_PORT=8000 \
    CONTRACTS_DIR=/app/contracts

COPY --from=api-build /out/api ./api
COPY --from=web-build /src/apps/web/.next/standalone ./
COPY --from=web-build /src/apps/web/.next/static ./apps/web/.next/static
COPY --from=web-build /src/apps/web/public ./apps/web/public
COPY contracts/ ./contracts/
COPY data/ ./data/
COPY scripts/start-container.mjs ./start-container.mjs
RUN chown -R node:node /app
USER node

EXPOSE 3000
HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=3 \
  CMD node -e "fetch('http://127.0.0.1:3000/api/v1/health').then(r => process.exit(r.ok ? 0 : 1)).catch(() => process.exit(1))"
CMD ["node", "/app/start-container.mjs"]
