# Experimenting with Go APIs

Messing about with building APIs in Go, how to structure things and whatnot.

I've created a bunch of folders, to experiment with different approaches:

- `std` - using the standard library for the HTTP server, `pgx` for database access
- `gin` - using the Gin framework for the HTTP server, `sqlc` + `pgx` for database access
- `fiber` - using the Fiber framework for the HTTP server, `gorm` for database access
- `echo` - using the Echo framework for the HTTP server, `gorm` for database access

Note that these are the major differences between the folders, but not the only ones. As I experiment with different approaches, other details will likely differ (e.g. validation, transaction script vs fp domain model vs oo domain model, etc).