# LG UltraFine Go Library [![CircleCI Build Status](https://circleci.com/gh/felicianotech/go-lguf.svg?style=shield)](https://circleci.com/gh/felicianotech/go-lguf) [![GitHub License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/felicianotech/go-lguf/master/LICENSE)

`go-lguf` is a Go library that interfaces with the LG UltraFine 4K monitor in order to adjust brightness from Linux.
This simple library is needed because this monitor was designed specifically for Apple computers and thus has no physical buttons.
Without the built-in features of macOS, adjusting brightness on this monitor wasn't possible.

Note: This was specifically made and tested for the UltraFine 4K 21.5" model 22MD4KA-B.
This may work with the newer 23" 2019 model, or even with the 5K model.
I don't know for sure, you can open an Issue or [tweet me](https://twitter.com/FelicianoTech) to let me know.


## Requirements

- **libusb** - Needed in order to talk to the monitor over USB.
Most modern Linux desktops will have this installed.
If not, on Ubuntu and Debian this can be installed by:

```bash
sudo apt install libusb-1.0
```
- *root* - In order to use `libusb`, root access is needed.


## Installation

`go-lguf` is a Go module and so the best way to use it is to import it into your own code and then run `go mod tidy` to get it downloaded.

```go
import(
	"github.com/felicianotech/go-lguf/lguf"
)
```


## Usage

The `lguf` package provides a `Connection` struct that does most of the work for you.
Just remember to `Close()` is.
Here's a simple program that simply prints the current brightness level and quits.

```go
package main

import (
	"fmt"
	"log"

	"github.com/felicianotech/go-lguf/lguf"
)

func main() {

	conn, err := lguf.NewConnection()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	value, err := conn.Brightness()
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("LG UltraFine 4K brightness is: %x", value)
}
```


## Development

This library is written and tested with Go v1.12+ in mind.
`go fmt` is your friend.
Please feel free to open Issues and PRs are you see fit.
Any PR that requires a good amount of work or is a significant change, it would be best to open an Issue to discuss the change first.


## License & Credits

This repository is licensed under the MIT license.
This repo's license can be found [here](./LICENSE).

This project would not exist if it wasn't for [Jean-François Beauchamp](https://github.com/velum)'s [repository](https://github.com/velum/lguf-brightness).
It's a C++ based, light-weight version of this.
