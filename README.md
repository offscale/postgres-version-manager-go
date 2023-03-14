postgres-version-manager-go
===========================
[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech)
[![License](https://img.shields.io/badge/license-Apache--2.0%20OR%20MIT%20OR%20CC0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Like [`rvm`](https://rvm.io)/[`nvm`](https://github.com/nvm-sh/nvm) but for PostgreSQL. Installs any of 138 different versions of PostgreSQL for your system.

## Development
Tested on Go 1.20, `git clone` this repo then `go build .`

## Usage

    NAME:
    postgres-version-manager-go - A new cli application
    
    USAGE:
    postgres-version-manager-go [global options] command [command options] [arguments...]
    
    DESCRIPTION:
    PostgreSQL version manager
    
    COMMANDS:
    help, h  Shows a list of commands or help for one command
    
    GLOBAL OPTIONS:
    --install value  Specific postgres version (default: latest)
    --start value    Specific postgres version (default: latest)
    --stop value     Specific postgres version (default: latest)
    --ls-remote      list available versions (default: false)
    --help, -h       show help

## License

Licensed under any of:

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or <https://www.apache.org/licenses/LICENSE-2.0>)
- MIT license ([LICENSE-MIT](LICENSE-MIT) or <https://opensource.org/licenses/MIT>)
- CC0 license ([LICENSE-CC0](LICENSE-CC0) or <https://creativecommons.org/publicdomain/zero/1.0/legalcode>)

at your option.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall be
licensed as above, without any additional terms or conditions.
