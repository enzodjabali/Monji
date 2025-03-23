# Stage 1: Build the app
FROM node:18-alpine AS builder
WORKDIR /app

# Copy package files and install dependencies (devDependencies are installed here)
COPY package*.json ./
RUN npm install

# Copy all source files and build the app
COPY . .
RUN npm run build

# Stage 2: Run the app
FROM node:18-alpine
WORKDIR /app

# Copy the built app from the builder stage
COPY --from=builder /app ./

EXPOSE 3000
CMD ["node", "build"]