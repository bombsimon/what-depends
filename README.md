# What depends?

This is a CLI tool to see what your code depends on. Instead of just looking for
URLs in your `go.mod` this package will query those URLs to fetch som
interesting information.

## Status

Not really useful at all - yet.

POC/pre-alpha/lab/WIP.

## Installation

```sh
go get -u github.com/bombsimon/what-depends/...
```

## Usage

```sh
$ what-depends
Dependencies for github.com/0x4b53/amqp-rpc
  ISC License
    * go-spew - Implements a deep pretty printer for Go data structures to aid in debugging

  BSD 3-Clause "New" or "Revised" License
    * uuid    - Go package for UUIDs based on RFC 4122 and DCE 1.1: Authentication and Security Services.

  MIT License
    * pretty  - Pretty printing for Go values
    * testify - A toolkit with common assertions and mocks that plays nicely with the standard library

  BSD 2-Clause "Simplified" License
    * amqp    - Go client for AMQP 0.9.1
```
