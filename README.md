# comdirect

A golang comdirect API Client implementation

## Library Usage

For an examplary usage see the [`e2e-test`](./cmd/e2e/command.go) file.

## Local usage

### Configuration

The configuration is done via a configuration file. The configuration file is expected to be in the current working directory and named `config.yaml`.

It can be stored in the following locations:
- `./config.yaml`
- `/etc/comdirect/config.yaml`
- `~/.comdirect/config.yaml`

For an example configuration see the [`config.yaml`](./config.example.yaml) file.