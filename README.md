# oVirt imageio

A port of [ovirt-imageio](https://github.com/oVirt/ovirt-imageio) to Go.

## Packages

Most of the daemon packages are implemented.

- images - images web server
- auth - authrization for images operations
- fileio - perform I/O to local file (file or block device)
- directio - helpers for doing direct I/O
- testutil - utilities for testing
- uuid - generates uuids version 4
- bench - benchmarks tools

## Testing

```
go test ./...
```

## Benchmarking

```
go build bench/recvfile.go
touch image
time ./recvfile -progress image 1024 </dev/zero
```

## Notes

- Support only aligned images and offsets. The Python version try to
  support unaligned images and offsets, make the code too complicated,
  and images should always be aligned.

- Ticket use {"mode": "rw"} instead of {"ops": ["read", "write"]}.
