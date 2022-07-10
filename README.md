# donut-fetch

A small utility for [donutdns](https://github.com/shoenig/donutdns) to download, cache, and create a master allow file from entries in a given sources.json file.

## Requirements

- Golang v1.18+

## Installation

```shell
go install github.com/neflyte/donut-fetch/cmd/donut-fetch@latest
```

## Usage

```
Usage:
  donut-fetch <sources.json> [flags]

Flags:
  -h, --help            help for donut-fetch
      --output string   output file; output to console if not specified
      --timeout uint    connection timeout in seconds (default 5)
```

## Notes

Downloaded host lists are cached in `${XDG_CACHE_HOME}/donut-fetch`.
`donut-fetch` state data is stored in `${XDG_CONFIG_HOME}/donut-fetch`.

Before downloading each host list from `sources.json`, an HTTP `HEAD` request is made to determine if the list has changed since it was last retrieved.
The `ETag` and `Last-Modified` headers from the result are checked, in that order.

## Building

```shell
make
```

### Clearing the cache

```shell
make clear-cache
```

### Clearing the state data

```shell
make clear-state
```

License: MIT
