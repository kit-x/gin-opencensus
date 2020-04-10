# gin-opencensus
[![Build Status](https://travis-ci.org/kit-x/gin-opencensus.svg?branch=master)](https://travis-ci.org/kit-x/gin-opencensus) [![GoDoc](https://godoc.org/github.com/kit-x/gin-opencensus?status.svg)](https://godoc.org/github.com/kit-x/gin-opencensus) [![Go Report Card](https://goreportcard.com/badge/github.com/kit-x/gin-opencensus)](https://goreportcard.com/report/github.com/kit-x/gin-opencensus) [![codecov](https://codecov.io/gh/kit-x/gin-opencensus/branch/master/graph/badge.svg)](https://codecov.io/gh/kit-x/gin-opencensus)  
opencensus middleware for gin

## Usage
```go
import "github.com/kit-x/gin-opencensus/ocgin"

e := gin.Default()

e.Use(ocgin.HandlerFunc())

```