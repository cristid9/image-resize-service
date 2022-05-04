
## Build & Run Server locally
```
go run .
```

## Run a sample request against the server

```bash
curl -X POST -H "Content-Type: application/json" -d @req.json 'http://localhost:8080/v1/resize?async=true'curl -X POST -H "Content-Type: application/json" -d @req.json 'http://localhost:8080/v1/resize?async=true'
```

Now in your browser, you can check one of the returned urls!
