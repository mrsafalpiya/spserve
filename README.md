# `spserve` - Serve files to current network with ease

`$ python3 -m http.server` but in golang and with better UI.

## Features

- No external dependencies.
- Excellent support for terminal browsers.
- JavaScript support is optional.
- Fast and easy to use.

## Installation

Assuming [proper golang setup](https://go.dev/doc/install), simply clone this
repo, cd to repo's current directory and run the following command:

```console
$ go install
```

## Usage

```console
$ spserve -h
Usage: spserve [options] root_dir

options:
  -h    Print help message
  -port int
        Port for the server to listen (default 8080)
```

```console
$ spserve .
2023/01/23 23:20:07 Serving "/home/safal" in 192.168.1.87:8080
```

```console
$ spserve -port 8008 ~/dl
2023/01/23 23:20:34 Serving "/home/safal/dl" in 192.168.1.87:8008
```

## Demo

![Demo](_docs/demo.png)

## License

GPLv3. See [COPYING](COPYING).
