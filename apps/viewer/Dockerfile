FROM node:22-alpine AS builder
RUN npm install -g pnpm@9.15.0
WORKDIR /app

# Copy workspace config for dependency install
COPY pnpm-lock.yaml pnpm-workspace.yaml package.json ./
COPY packages/api-types/ packages/api-types/
COPY packages/api-client/ packages/api-client/
COPY packages/tsconfig/ packages/tsconfig/
COPY packages/tailwind-config/ packages/tailwind-config/
COPY apps/viewer/ apps/viewer/

# Install all workspace dependencies
RUN pnpm install --frozen-lockfile

# Build only this app
RUN pnpm --filter @neshiman/viewer build

FROM nginx:alpine
COPY --from=builder /app/apps/viewer/dist /usr/share/nginx/html
COPY apps/viewer/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
