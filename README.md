# cryptik
<p align="center">
  <img src="assets/cryptik.png" height="200" alt="cryptik logo" />
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/shubhvish4495/cryptik.svg)](https://pkg.go.dev/github.com/shubhvish4495/cryptik)
[![Go Report Card](https://goreportcard.com/badge/github.com/shubhvish4495/cryptik)](https://goreportcard.com/report/github.com/shubhvish4495/cryptik)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
---



---

## Overview

**Cryptik** is a secure, flexible One-Time Password (OTP) generation and validation library for Go (Golang). It features built-in caching, cryptographic random generation, and easy integration for any Go application requiring OTP-based authentication.

---

## Features

* Cryptographically secure OTP generation
* Configurable OTP length (default: 6 digits)
* In-memory cache with expiration support
* Thread-safe operations
* Automatic cleanup of expired OTPs
* Easy API for developers

---

## Installation

```bash
go get github.com/shubhvish4495/cryptik
```

---

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/shubhvish4495/cryptik"
)

func main() {
    otpService, err := cryptik.NewService(cryptik.CryptikConfig{
  		Length: 6,
  		Cache:  cache.GetCache(),
  	})
    if err != nil {
        panic(err)
    }

    secret := "user123"
    otp, err := otpService.GenerateOTP(secret)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Generated OTP: %s\n", otp)

    isValid, err := otpService.ValidateOTP(secret, otp)
    if err != nil {
        panic(err)
    }
    fmt.Printf("OTP is valid: %v\n", isValid)
}
```

---

## Detailed Features

### OTP Generation

* Uses `crypto/rand` for secure randomness
* OTP length configurable (default: 6 digits)
* Automatically stores OTP in internal cache with expiration
* Default expiration: 10 minutes

### OTP Validation

* Compares user input with cached OTP
* Removes OTP after successful validation (single-use)
* Thread-safe

### Built-in Cache

* In-memory store
* Auto-removal of expired OTPs
* Fully thread-safe
* Supports custom implementation

---

## Custom Cache

Implement the `Cache` interface:

```go
type Cache interface {
    Get(key string) (any, bool)
    Set(key string, value any, expiration int64) error
    Delete(key string)
    Exists(key string) bool
    Clear()
    RemoveExpiredEntries()
}
```

Use it like this:

```go
otpService, err := cryptik.NewService(cryptik.cryptikServiceConfig{
    Cache:  myCustomCache,
    Length: 6,
})
```

---

## Configuration

```go
type cryptikServiceConfig struct {
    Cache  cache.Cache // optional
    Length int         // optional, default: 6
}
```

---

## Error Handling

```go
var ErrInvalidOTP = errors.New("invalid OTP provided")
```

---

## Thread Safety

All major operations (generate, validate, cache access) are safe for concurrent use.

---

## Best Practices

1. Use a unique identifier for each OTP (e.g., email, userID)
2. Keep OTP lifespan short (≤10 minutes)
3. Never allow reuse of OTPs
4. Use at least 6 digits for strong security

---

## License

MIT License. See `LICENSE` file for details.

---

## Contributing

Issues and PRs welcome! Help make Cryptik better.

---

## Keywords (SEO)

`golang`, `otp`, `go otp`, `secure otp`, `csprng`, `go library`, `authentication`, `go module`, `otp generation`, `one-time password`
