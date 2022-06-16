## Notes

Docker build locally
```
docker build -t auth .
```

Docker run locally with attached STDOUT etc
```
docker run -e HOST_ADDRESS=:3000 -p 3000:3000 auth
```

Docker run locally detached
```
docker run -d -e HOST_ADDRESS=:3000 -p 3000:3000 auth
```
