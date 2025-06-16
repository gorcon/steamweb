# steamweb
[![GitHub Build](https://github.com/gorcon/steamweb/workflows/build/badge.svg)](https://github.com/gorcon/steamweb/actions)
[![Go Coverage](https://github.com/gorcon/steamweb/wiki/coverage.svg)](https://raw.githack.com/wiki/gorcon/steamweb/coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/gorcon/steamweb)](https://goreportcard.com/report/github.com/gorcon/steamweb)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/gorcon/steamweb)

Steam Web API Implementation in Golang.

## API Specifications
Steam API described in the [valve documentation](https://developer.valvesoftware.com/wiki/Steam_Web_API).

## Install
```text
go get github.com/gorcon/steamweb
```

See [Changelog](CHANGELOG.md) for release details.

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	steamweb "github.com/gorcon/steamweb/steamwebdraft"
)

func main() {
	client := steamweb.NewClient(&steamweb.Config{Key: "{Steam API Key}"})

	servers, err := client.GetServerList(&steamweb.GetServerListFilter{}) // Set filters here
	if err != nil {
		log.Fatal(err)
	}

	js, _ := json.Marshal(servers)

	fmt.Println(string(js))
}
```

## Requirements
Go 1.23 or higher

## Contribute
Contributions are more than welcome!

If you think that you have found a bug, create an issue and publish the minimum amount of code triggering the bug, so
it can be reproduced.

If you want to fix the bug then you can create a pull request. If possible, write a test that will cover this bug.

## License
MIT License, see [LICENSE](LICENSE)
