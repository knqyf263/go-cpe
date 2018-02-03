# go-cpe

[![Build Status](https://travis-ci.org/knqyf263/go-cpe.svg?branch=master)](https://travis-ci.org/knqyf263/go-cpe)
[![Coverage Status](https://coveralls.io/repos/github/knqyf263/go-cpe/badge.svg?branch=initial)](https://coveralls.io/github/knqyf263/go-cpe?branch=initial)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](/LICENSE)

A Go library for [CPE(A Common Platform Enumeration 2.3)](https://cpe.mitre.org/specification/)

`go-cpe` is a implementation of the CPE Naming and Matching algorithms, as described in NIST IRs [7695](https://csrc.nist.gov/publications/detail/nistir/7695/final) and [7696](https://csrc.nist.gov/publications/detail/nistir/7696/final).  

For the reference Java implementation, see: https://cpe.mitre.org/specification/CPE_Reference_Implementation_Sep2011.zipx

# Installation and Usage

Installation can be done with a normal go get:

```
$ go get github.com/knqyf263/go-cpe
```

## Compare
See [example](/example)

```
package main

import (
	"fmt"

	"github.com/knqyf263/go-cpe/matching"
	"github.com/knqyf263/go-cpe/naming"
)

func main() {
	wfn, err := naming.UnbindURI("cpe:/a:microsoft:internet_explorer%01%01%01%01:8%02:beta")
	if err != nil {
		fmt.Println(err)
		return
	}
	wfn2, err := naming.UnbindFS(`cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:*:*:*:*:*`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(matching.IsDisjoint(wfn, wfn2))
	fmt.Println(matching.IsEqual(wfn, wfn2))
	fmt.Println(matching.IsSubset(wfn, wfn2))
	fmt.Println(matching.IsSuperset(wfn, wfn2))
}
```

# Contribute

1. fork a repository: github.com/knqyf263/go-cpe to github.com/you/repo
2. get original code: `go get github.com/knqyf263/go-cpe`
3. work on original code
4. add remote to your repo: git remote add myfork https://github.com/you/repo.git
5. push your changes: git push myfork
6. create a new Pull Request

- see [GitHub and Go: forking, pull requests, and go-getting](http://blog.campoy.cat/2014/03/github-and-go-forking-pull-requests-and.html)

----

# License
MIT

# Author
Teppei Fukuda