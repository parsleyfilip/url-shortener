package main

import (
	"net/http"
)

// String constants for HTML content
const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>URL Shortener</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 30px;
            margin-top: 50px;
        }
        h1 {
            color: #2c3e50;
            margin-top: 0;
            text-align: center;
        }
        form {
            display: flex;
            flex-direction: column;
        }
        label {
            font-weight: bold;
            margin-bottom: 5px;
        }
        input[type="url"] {
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
            margin-bottom: 20px;
        }
        button {
            background-color: #3498db;
            color: white;
            border: none;
            padding: 12px 20px;
            border-radius: 4px;
            font-size: 16px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        button:hover {
            background-color: #2980b9;
        }
        #result {
            margin-top: 20px;
            padding: 15px;
            border-radius: 4px;
            background-color: #f8f9fa;
            display: none;
        }
        .short-url {
            font-weight: bold;
            color: #3498db;
            word-break: break-all;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            color: #777;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>URL Shortener</h1>
        <form id="shorten-form">
            <label for="url-input">Enter a long URL to make it shorter:</label>
            <input type="url" id="url-input" placeholder="https://example.com/very/long/url" required>
            <button type="submit">Shorten URL</button>
        </form>
        <div id="result">
            <p>Your shortened URL:</p>
            <p class="short-url" id="short-url"></p>
        </div>
    </div>
    <div class="footer">
        <p>Simple URL Shortener Service</p>
    </div>

    <script>
        document.getElementById('shorten-form').addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const urlInput = document.getElementById('url-input');
            const resultDiv = document.getElementById('result');
            const shortUrlElement = document.getElementById('short-url');
            
            try {
                const response = await fetch('/shorten', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        url: urlInput.value
                    })
                });
                
                if (!response.ok) {
                    throw new Error('Failed to shorten URL');
                }
                
                const data = await response.json();
                
                shortUrlElement.textContent = data.short_url;
                shortUrlElement.href = data.short_url;
                resultDiv.style.display = 'block';
                
            } catch (error) {
                alert('Error: ' + error.message);
            }
        });
    </script>
</body>
</html>
`

// ServeIndex serves the HTML homepage
func serveIndex(w http.ResponseWriter, r *http.Request) {
	// Only serve the index page on the root path with GET method
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		return
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexHTML))
} 