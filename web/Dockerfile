FROM hub.megan.ir/oven/bun:latest AS base
WORKDIR /app

FROM base AS deps
COPY package.json ./
COPY bunfig.toml ./
RUN bun install

FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV DATABASE_URL=postgresql://postgres:postgres@db:5432/hotel
RUN bun run build

FROM base AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV DATABASE_URL=postgresql://postgres:postgres@db:5432/hotel

COPY --from=builder /app/.output ./.output
COPY package.json ./

EXPOSE 3000
CMD ["bun", "run", ".output/server/index.mjs"]
