[![Go](https://github.com/fairking/di-plus/actions/workflows/go.yml/badge.svg)](https://github.com/fairking/di-plus/actions/workflows/go.yml)

# Dependency Injection

Everything you need for [Dependency Injection (DI)](https://wikipedia.org/wiki/Dependency_injection) in golang.

## Features
This package provides the following features:
- Register services with a single instance (singleton only) or with a function (factory);
- Use of generic function to get your desired service;
- Three service lifetimes: singleton, scoped, transient;
- Async initialization (coming soon)

## Service Lifetimes
A service can have the following lifetimes:
- `Transient` - a new instance is created every time it is requested
- `Singleton` - a single, new instance is created the first time it is requested
- `Scoped` - a new instance is created once per ServiceScope the first time it is requested

## Scope Scenarios
There scenarios where a service needs to be scoped; for example, for the lifetime of a HTTP request. A service definitely shouldn't live for the life of the application (e.g. `singleton`), but it also shouldn't be created each time it's requested within the request (e.g. `transient`). A scoped service lives for the lifetime of the container (`ServiceScope`) it was created from.

## Installation
```bash
go get github.com/fairking/di-plus
```

## Usage
```golang
package main

import (
    "github.com/fairking/di-plus"
)

type MySettingsService struct {
	Connection string
}

type MyDatabaseService struct {
	Settings MyDatabaseService
}

// Register Services
services := NewServiceProvider().
    // Singleton
    Register(
        GetTypeOf[MySettingsService](),
        ServiceTypeEnum.Singleton,
        func(s IServiceProvider) (interface{}, error) {
            return MySettingsService{Connection: "server=127.0.0.1;uid=root;pwd=12345;database=test"}, nil
        },
    ).
    // Scoped
    Register(
        GetTypeOf[MyDatabaseService](),
        ServiceTypeEnum.Scoped,
        func(s IServiceProvider) (interface{}, error) {
            return MyDatabaseService{Settings: GetService[*MySettingsService](s)}, nil
        },
    ).
    Build()

// Get Singleton
settings := services.GetService("MySettingsService").(*MySettingsService) // Creates a new settings value
// or
settings2 := GetService[*MySettingsService](services) // Returns the existing settings value
// settings and settings2 pointing to the same value

// Get Scope
{
    scope := services.GetScope()
    db := GetService[*MyDatabaseService](scope) // Creates a new database service value
    // or
    db_, err := GetServiceOr[*MyDatabaseService](scope) // Returns the existing database service value

    settings3 := GetService[*MySettingsService](services) // Returns the existing settings value
}

// Get Scope
{
    scope := services.GetScope()
    db2 := GetService[*MyDatabaseService](scope) // Creates a new database service value
    // or
    db2_, err := GetServiceOr[*MyDatabaseService](scope) // Returns the existing database service value

    settings3 := GetService[*MySettingsService](services) // Returns the existing settings value
}

// db and db_ is the same value, but db and db2 are different

```

## Contribution
If you have questions please create a new [issue](https://github.com/fairking/di-plus/issues/new/choose).

## Donations

Donate with [Ӿ nano crypto (XNO)](https://nano.org).

[![di_plus Donations](https://gitlab.com/fairking/sqlkata.queryman/-/raw/master/Resources/Donations_QRCode_nano_1sygjbke.png)](https://nanocrawler.cc/explorer/account/nano_1sygjbkepdcu5diiekf15ar6m6utfgf9rr9tkd6zi8mkq7yza34kiyjpgt9g)

Thank you!