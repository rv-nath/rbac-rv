# RBAC Library

This is a Role-Based Access Control (RBAC) library for Go applications, designed to support role and permission checks at a granular level. It supports scope-based permissions such as `self`, `group`, and `global` and is structured to be easily integrated into an API layer, including use with API gateways like [KrakenD](https://www.krakend.io/).

## Features

- Flexible role and permission management
- Scope-based checks for self, group, and global access
- Dynamic resource and intent parsing
- Easy integration with API gateways and other Go applications

## Installation

To install the module, use `go get`:

```bash
go get github.com/your-username/your-repo@v1.0.0

