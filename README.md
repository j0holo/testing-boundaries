# Example of testing boundaries in an application

This is an example where you could focus your tests on boundaries in your application. A boundary is contract or interface between two components. I go more in depth on the subject in my blog: https://blog.joholo.com/bdd-and-boundaries-part-2/.

The components in this application are:

```
HTTP server -> business logic -> database layer
                             \-> queue layer <- queue worker
```

By focussing our tests on the HTTP server and queue worker our implementation of the database and queue layer is more
flexible. Adding a storage layer for our business logic doesn't break our tests. If we upgrade our third party S3
library which changes the implementation of the storage layer is doesn't break our tests.

When the output of the HTTP server is difficult to test you might focus your tests on the business layer instead to keep
testing easier. You could write a simple test against your HTTP server if possible to see if a request with status 200
is possible. Just to give your more confidence.

## Run the tests

```bash
docker-compose up -d
go test -v -cover ./...
```

## Current test coverage

| filename | coverage % |
|---|---|
| business.go | 92.3% |
| database.go | 100% |
| main.go | 0% |
| routes.go | 78.1% |
| worker.go | 85.7% |

With a total of 71.2% which is mainly because `main.go` doesn't have any tests. This can be improved by adding functions
for Postgres and Redis setup that can be reused in the tests. Without `main.go` the coverage would be 89.0%.
