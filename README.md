# URL Shortener

A simple URL shortener service built with Go and MongoDB.
I made this for a friend!

## Features

- Shorten long URLs into compact, shareable links
- Redirect shortened URLs to their original destinations
- Track visit counts for each shortened URL
- RESTful API for URL shortening
- MongoDB for fast and scalable storage

## Requirements

- Go 1.17 or higher
- MongoDB server or cloud instance

## Environment Variables

The application uses the following environment variables:

- `MONGODB_URI`: MongoDB connection string (default: `mongodb://localhost:27017`)
- `MONGODB_DATABASE`: MongoDB database name (default: `url_shortener`)
- `PORT`: HTTP server port (default: `8080`)

## Setup Instructions for Strato Webhosting

1. Upload the compiled Go binary to your Strato webhosting account.
2. Configure the MongoDB connection in your Strato environment.
3. Make sure the binary has execution permissions.
4. Create a startup script or use Strato's application management to start the service.

## Building from Source

```bash
# Clone the repository
git clone https://github.com/101179/url-shortener.git
cd url-shortener

# Download dependencies
go mod download

# Build the application
go build -o url-shortener .
```

## Usage

### Starting the Service

```bash
# Run with default settings
./url-shortener

# Run with custom settings
PORT=9000 MONGODB_URI="mongodb+srv://user:password@cluster.mongodb.net" MONGODB_DATABASE="my_urls" ./url-shortener
```

### API Endpoints

#### Shorten a URL
```
POST /shorten
Content-Type: application/json

{
  "url": "https://example.com/very/long/url/that/needs/shortening"
}
```

Response:
```json
{
  "original_url": "https://example.com/very/long/url/that/needs/shortening",
  "short_url": "http://your-domain.com/Ab3Def9g"
}
```

#### Redirect to Original URL
```
GET /{shortCode}
```

This will redirect the user to the original URL.

#### Health Check
```
GET /health
```

Returns "OK" if the service is running.

## License

MIT
