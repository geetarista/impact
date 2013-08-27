# Go Impact ![Build Status](https://travis-ci.org/geetarista/impact.png?branch=master)

This is a server written in [Go](http://golang.org) for the [ImpactJS](http://impactjs.com/) game engine.

So far, the only thing it does is replace the PHP server that comes with ImpactJS.

## Installation and Usage

###. You do not have Go installed

You can just download the [binary version](https://raw.github.com/geetarista/impact/master/impact).

Then make it executable for your user:

```bash
chmod +x impact
```

Then copy it to the root of your game's path and run it:

```bash
./impact
```

### You have Go installed

Run the following to install the server:

```bash
go get github.com/geetarista/impact
```

Then just start the server:

```bash
impact
```

## Server Port

The default port the server will listen on is 8080, but you can override that like so:

```bash
PORT=8888 impact
```

## Tests

To run the tests, you just need Go installed on your system and just run:

```bash
go test
```

## Disclaimer

This server has only been tested on Mac OS X for now, so if you discover issues on other platforms, please [create an issue](https://github.com/geetarista/impact/issues) so I can address them.

## License

MIT. See `LICENSE`.
