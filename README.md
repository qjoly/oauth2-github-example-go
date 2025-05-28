# GitHub OAuth2 Authentication in Go

A simple Go application demonstrating GitHub OAuth2 authentication flow.

## Setup

1. **Register a GitHub OAuth Application**:
   - Go to [GitHub Developer Settings](https://github.com/settings/developers)
   - Click "New OAuth App"
   - Set callback URL to `http://localhost:8080/auth/github/callback`

2. **Configure Environment Variables**:
   ```bash
   echo "GITHUB_CLIENT_ID=your_client_id" > .env
   echo "GITHUB_CLIENT_SECRET=your_client_secret" >> .env
   ```

3. **Install Dependencies**:
   ```bash
   go get golang.org/x/oauth2
   ```

## Running the Application

```bash
go run main.go
```

Visit [http://localhost:8080](http://localhost:8080) in your browser.


