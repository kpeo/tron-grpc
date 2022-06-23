# tron-grpc

Simple Go GRPC-client for the TRON blockchain.

Currently project uses [gotron-sdk](https://github.com/fbsobreira/gotron-sdk) for compiled protos from [tronprotocol](https://github.com/tronprotocol/protocol).

As far as gotron-sdk is not updated frequently one of todos is to compile tronprotocol directly.

## Install

```
  go build
```

## Usage

You can use this project as library that simplifies interaction with TRON nodes (currently for accounts, transfers and basic info).

### Setting environment

To get specific variables more secure way or for Docker you can fill the required environment variables in .env file or set them via standart environment variables.

## Testing

Currently to test the project you need to fill testing variables in `tron/*_test.go`.

## TODO

- Full test coverage
- Direct tronprotocol compilation
- Extension of the functional
- Adding CLI and Web-interfaces

