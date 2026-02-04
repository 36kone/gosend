# GoSend

**GoSend** is a command-line tool written in Go to send HTTP requests in an interactive and colorful way.  
It works similarly to `curl`, but focuses on interactivity, colorized output, and ease of use for developers.

---

## 🏷 Features

- Supports **GET, POST, PUT, DELETE** HTTP methods  
- Interactive input for URL, method, body (JSON/Form), and headers  
- Add custom headers, including **Authorization tokens**, for authenticated requests  
- Colorized output for **JSON and status codes**  
- Lightweight and fast alternative for quick API testing 

---

## 💻 Installation (Linux)

You can eventually install via **.deb**, but for development:

```bash
# Clone the repository
git clone https://github.com/36kone/gosend.git
cd gosend

# Build the binary
go build -o gosend ./cmd

# Optional: move to /usr/local/bin to run from anywhere
sudo mv gosend /usr/local/bin/

# Then run:
gosend
```
# 🚀 Usage

---

## Interactive Mode
```bash
# Run GoSend in interactive mode
gosend

# You will be prompted to provide:

- URL

- Protocol

- HTTP Method

- Body (JSON or Form)

- Headers (optional)

Tip: You can use autocomplete while typing URLs or HTTP methods.
By default, the protocol is `https` and the method is `GET`.

```
## Flags Mode (future)
```bash
# Run GoSend in flag mode

# Simple GET request
gosend -u "https://api.example.com/health" -X GET

# GET request with custom headers (e.g., Authorization Bearer token)
gosend -u "https://api.example.com/secure-data" -X GET \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json"
```
🎨 Colorized Output
JSON: Green for strings, Blue for keys/structures, Yellow for :

```bash
# Status Code:

1xx → Gray

2xx → Green

3xx → Yellow

4xx → Light Red

5xx → Strong Red

```
🤝 Contributing
Pull requests are welcome!
Feel free to open issues or suggest improvements.

📝 License
MIT License © Caio Herrera
