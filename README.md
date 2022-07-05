# Testing CORS with gRPC Gateway

1. Start the gRPC Gateway

```bash
go run cmd/gateway/main.go
```

2. Go to `google.com` and execute request from browser console

```javascript
fetch(
  'http://localhost:8081',
  {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' }
  }
).then(resp => resp.text()).then(console.log)
```
