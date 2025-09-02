# Experimenting with Go APIs

Messing about with building APIs in Go, how to structure things and whatnot.

I've created a bunch of folders, to experiment with different approaches:

- `std` - using the standard library for the HTTP server, `pgx` for database access
- `gin` - using the Gin framework for the HTTP server, `sqlc` + `pgx` for database access
- `fiber` - using the Fiber framework for the HTTP server, `gorm` for database access
- `echo` - using the Echo framework for the HTTP server, `gorm` for database access
- `chi` - using the Chi framework for the HTTP server, `TBD` for database access

Note that these are the major differences between the folders, but not the only ones. As I experiment with different approaches, other details will likely differ (e.g. validation, transaction script vs fp domain model vs oo domain model, etc).

Some inspiration for the structure and features came from:
- [How I write HTTP services in Go after 13 years](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/)
- [How I've been building APIs and microservices lately (feat. C# & .NET)](https://blog.codingmilitia.com/2025/06/11/how-ive-been-building-apis-and-microservices-lately-feat-csharp-dotnet/)

## Some notes taken along the way

### Web frameworks

- `std`
  - It's surprisingly straightforward to build an API using just the standard library
  - Maybe a bit more boilerplate code than others, but it's not too bad, creating an internal boilerplate package with some helpers should do the trick
- `gin`
  - TODO 👷‍♀️
- `fiber`
  - Less boilerplate when compared to `std`
  - The returning of errors from handlers is nice, leaving to the error handling middleware to centralize returning the error response
  - The `ctx` situation is a bit confusing, as we can get `c.Context()` or `c.UserContext()`, the first one being from fasthttp, with request information, while the latter is a standard `context.Context`, where stuff like OpenTelemetry span info is stored. Plus, neither seems to be cancelled when the client disconnects, which is something the standard library support.
- `echo`
    - TODO 👷‍♀️
- `chi`
    - TODO 👷‍♀️

### Database access
- `pgx`
  - Overall easy to use
- `sqlc`
  - TODO 👷‍♀️
- `gorm`
  - TODO 👷‍♀️