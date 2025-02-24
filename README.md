# comdirect

A golang comdirect API Client implementation

## Library Usage

For an examplary usage see the [`e2e-test`](./cmd/e2e/command.go) file.

## Local usage

### Configuration

The configuration is done via a configuration file. The configuration file is expected to be in the current working directory and named `config.yaml`.

It can be stored in the following locations:
on unix:
- `./config.yaml`
- `/etc/comdirect/config.yaml`
- `~/.comdirect/config.yaml`

or on windows:
- `C:\\ProgramData\\comdirect\\config.yaml`
- `C:\\Users\\%USERNAME%\\AppData\\Roaming\\comdirect\\config.yaml`
- `.\\config.yaml`

For an example configuration see the [`config.yaml`](./config.example.yaml) file.

### Examples

Find further usage by running the following commands:

```bash
comdirect --help
```

#### Get account balances

```bash
comdirect account balances -o yaml
```

#### Get account transactions

```bash
comdirect account transactions <account_id>
```