#!/bin/bash
# Setup script for running URL Shortener on Strato Webhosting

# Create necessary directories
mkdir -p logs

# Download MongoDB if it's not available on Strato
# Note: This is just an example, you'll need to adapt this to Strato's environment
# If Strato doesn't allow MongoDB installation, consider using a cloud MongoDB service like MongoDB Atlas

# Setting environment variables - adjust these as needed
export PORT=8080
export MONGODB_URI="mongodb://localhost:27017"  # Or your cloud MongoDB URI
export MONGODB_DATABASE="url_shortener"

# Make the Go binary executable
chmod +x url-shortener

# Start the application in the background and log output
nohup ./url-shortener > logs/app.log 2>&1 &
echo $! > app.pid

echo "URL Shortener started on port $PORT"
echo "Process ID saved in app.pid"
echo "Logs available in logs/app.log" 