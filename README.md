rwc [![License](http://img.shields.io/:license-gpl3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0.html) [![GoDoc](http://godoc.org/github.com/opennota/rwc?status.svg)](http://godoc.org/github.com/opennota/rwc) [![Build Status](https://travis-ci.org/opennota/rwc.png?branch=master)](https://travis-ci.org/opennota/rwc)
===

Package rwc provides a pseudo-Russian word constructor.

## Install

    go get -u github.com/opennota/rwc

## Use

``` Go
package main
import (
	"fmt"
	"github.com/opennota/rwc"
)
func main() {
	fmt.Println(rwc.Word(7))
	fmt.Println(rwc.WordMask("зVC...я"))
}
```

## See also

http://kirsanov.com/rwc/

