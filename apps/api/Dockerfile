##############################
# Stage 1: Build the web app
##############################
FROM node:18-alpine AS web_builder
WORKDIR /app
# Copy package files and install dependencies
COPY apps/web/package*.json ./
RUN npm install
# Copy the entire web source code and build the app
COPY apps/web/ .
RUN npm run build
# (Optional) List the directory for debugging to see what’s produced
RUN ls -la /app

##############################
# Stage 2: Build the Go API
##############################
FROM golang:1.18-alpine AS api_builder
# Install build dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /api
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download
COPY apps/api/ .
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api

##############################
# Stage 3: Final Runtime Image
##############################
FROM node:18-alpine
# Remove default entrypoint by overriding it with an empty one.
ENTRYPOINT []
# Install runtime libraries for SQLite
RUN apk add --no-cache sqlite-libs
WORKDIR /app

# Copy the compiled Go API binary
COPY --from=api_builder /api/main ./api

# Copy the entire web build (this copies everything from the web builder stage)
COPY --from=web_builder /app /web

# (Optional) In case your preview command needs production dependencies.
WORKDIR /web
RUN npm install --omit=dev

# Create a startup script at the root
WORKDIR /
RUN echo '#!/bin/sh\n\
echo "Starting Go API..."\n\
/app/api &\n\
echo "Starting Web Preview..."\n\
cd /web && npm run preview -- --host 0.0.0.0 --port 3000\n\
wait' > /start.sh && chmod +x /start.sh

# Debug: list the root to confirm /start.sh exists
RUN ls -la /

# Expose API on 8080 and Web on 3000
EXPOSE 8080 3000

# Use an explicit shell command to run the startup script.
CMD ["sh", "/start.sh"]