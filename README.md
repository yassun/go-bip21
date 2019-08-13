# go-bip21

[![Build Status](https://secure.travis-ci.org/yassun/go-bip21.png?branch=master)](http://travis-ci.org/yassun/go-bip21)
[![Coverage Status](https://coveralls.io/repos/yassun/go-bip21/badge.svg?branch=master)](https://coveralls.io/r/yassun/go-bip21?branch=master)
[![GoDoc](https://godoc.org/github.com/yassun/go-bip21?status.svg)](https://godoc.org/github.com/yassun/go-bip21)
[![license](https://img.shields.io/badge/license-MIT-4183c4.svg)](https://github.com/yassun/go-bip21/blob/master/LICENSE)

go-bip21 is an open source library to handle the URI based on the [BIP-21](https://github.com/bitcoin/bips/blob/master/bip-0021.mediawiki) standard.
 
# Install

```bash
$ go get github.com/yassun/go-bip21
```

# Usage

Parse the URI `bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=20.3&label=Luke-Jr`.

```Go
u, err := bip21.Parse("bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=20.3&label=Luke-Jr")
if err != nil {...}

// &{UrnScheme:bitcoin Address:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W Amount:20.3 Label:Luke-Jr Message: Params:map[]}
fmt.Printf("%+v\n", u)
```

Build the URI `bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=20.3&label=Luke-Jr`

```Go
u := &bip21.URIResources{
  UrnScheme: "bitcoin",
  Address:   "175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W",
  Amount:    20.3,
  Label:     "Luke-Jr",
  Message:   "",
  Params:    make(map[string]string),
}

uri, err := u.BuildURI()
if err != nil {...}

// bitcoin:175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W?amount=20.3&label=Luke-Jr
fmt.Printf("%+v\n", uri)
```
