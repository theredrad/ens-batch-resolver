# ENS batch domain resolver (.eth domain)
A simple program to check a batch of ENS domains availability.

## Configure
Configs store in `config.yaml` file next to the main file. an example of the config file exists.

### `client-endpoint`
To connect to the Ethereum network, you need to create an account in https://infura.io & put your API URL in the config.yml as `client-endpoint`.

### `list-file`
The `list-file` is the path of a json file with an object with an array of domains which is named `domains`.
Example:
```json
{
  "domains": [
    "abacus",
    "abased",
    "abated",
    "abates",
    "abayas"
  ]
}
```
### `output-file`
The `output-file` is the path of output file.

## Run
You can pull the source & run it via `go run main.go` or download the binaries from the release page. make sure the binary file is executable & run it from the terminal `./resolver`.
