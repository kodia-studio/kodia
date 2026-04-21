# Builder Stage
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files for better caching
COPY package*.json ./

# Install all dependencies (including devDependencies for build)
RUN npm ci

# Copy the rest of the application
COPY . .

# Build the application
RUN npm run build

# Remove development dependencies
RUN npm prune --production

# Run Stage
FROM node:22-alpine

WORKDIR /app

# Install runtime dependencies (like curl for healthchecks)
RUN apk add --no-cache curl

# Copy necessary files from builder
COPY --from=builder /app/package.json ./
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules

# Environment variables
ENV NODE_ENV=production
ENV PORT=3000

# Expose the application port
EXPOSE 3000

# Start the application
CMD ["node", "build"]
