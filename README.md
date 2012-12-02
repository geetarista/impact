# Go Impact

This is a server written in [Go](http://golang.org) for the [ImpactJS](http://impactjs.com/) game engine.

So far, the only thing it does is replace the PHP server that comes with ImpactJS.

## Installation and Usage

1. You have Go installed

Run the following to install the server:

```bash
go get github.com/geetarista/impact
```

Then just start the server:

```bash
impact
```

2. You do not have Go installed

You can just download the [binary version](https://raw.github.com/geetarista/impact/master/impact).

Then make it executable for your user:

```bash
chmod +x impact
```

Then copy it to the root of your game's path and run it:

```bash
./impact
```

## License

MIT. See `LICENSE`.
