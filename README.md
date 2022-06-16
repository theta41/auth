## Notes

### Go run
Run
```
go run cmd/auth/main.go -c config.yaml
```

### Docker

Build locally
```
docker build -t auth .
```

Run locally with attached STDOUT etc
```
docker run -p 3000:3000 auth
```

Run locally detached
```
docker run -d -p 3000:3000 auth
```
