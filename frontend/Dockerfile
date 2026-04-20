# ── Stage 1: Build ──
FROM node:22-alpine AS builder

WORKDIR /build

# Cache dependency installs
COPY package.json package-lock.json* ./
RUN npm install

# Build Nuxt application
COPY . .
RUN npm run build


FROM builder AS dev

WORKDIR /build

CMD ["npm", "run", "dev"]

# ── Stage 2: Run ──
FROM node:22-alpine

WORKDIR /app

# Copy only the production output
COPY --from=builder /build/.output .output

ENV HOST=0.0.0.0
ENV PORT=3000
ENV NODE_ENV=production

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
