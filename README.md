# go-pgp-tool

## Usage
build code
```bash
$ go build -o go-pgp-tool main.go
```
show usage
```bash
$ ./go-pgp-tool -h
```

## test
```bash
$ go test -v ./...
```

### options
```
  Usage options:
    -h show usage
    -d decrypt mode
    -e encrypt mode
    -y force override output file
    -p password
      private key password with decrypt usage
    -i file path
      source file path
    -o [file path|-]
      output file path or stdout
    -k file path
      key file path
```