# `pfin` - Personal Finance Tracker

Very simple and personal finance tracker, please do not use :)

Data is written to `(data_dir)/pfin.jsonl` configurable through `~/.config/pfin/config.toml`. A good default will be written out if it doesnt exist.

Example:

```sh
# help
$ go run ./cmd/main.go -h

# add
$ go run ./cmd/main.go -testmode add -etype expense -amount 1000 -category other -note 'payment'
$ go run ./cmd/main.go -testmode add -etype expense -amount 1000 -category other -note 'payment' -date "2024-08-11"

# -testmode overrides config to be testdata/testconfig.toml and writes to testdata/
$ go run ./cmd/main.go -testmode add -etype expense -amount 1000 -category other -note 'payment'
$ go run ./cmd/main.go -testmode add -etype expense -amount 1000 -category other -note 'payment' -date "2024-08-11"

# specify a custom config other than ~/.config/pfin/config.toml
$ go run ./cmd/main.go -config [path_to_custom_config] [...args]
```

There are some tests. Run them using `go test`

```sh
$ go test ./... # run all tests
```
