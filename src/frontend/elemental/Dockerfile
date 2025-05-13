# Use Node.js image
FROM node:18

# Set working directory
WORKDIR /app

# Copy package.json and install dependencies
COPY package*.json ./
RUN npm install

# Copy the rest of the frontend app
COPY . .

# Expose port (default for Vite, Next.js, etc.)
EXPOSE 3000

# Run the dev server
CMD ["npm", "run", "dev"]
