# PortDomainService
The PortDomainService provides an HTTP endpoint that accepts a URI of a file with port data in JSON to be imported into
a database. Due to the potentially long processing time for large files, the response to the client is sent back before
the data is processed - all work is done asynchronously.

Endpoint is available at `/ports`, and accepts data as a JSON object:

```
curl --request POST 'localhost:8088/ports' \
--data-raw '{
    "file_uri":"https://www.example.com/ports.json"
}'
```

Launch:
- add .env file, you can start by copying `.env.example`
- `docker-compose build`
- `docker-compose up`

Test:
```
go test ./... -race
```

TODO:
- metrics
- better error handling, continue processing files on some types of errors, retry on others
- use gRPC, possibly allow replaying messages on some types of errors
- complete PostgreSQL DB setup (docker-compose + service)
