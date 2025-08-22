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