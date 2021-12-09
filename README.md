# Site24x7 Command Line Client

A CLI tool for interacting with the [Site24x7 API](https://www.site24x7.com/help/api/#introduction).

This is currently _unreleased software_. The content of this README is anticipatory more than actual.

## Prerequisites

1. A Site24x7 account
2. A registered [Site24x7 application](https://api-console.zoho.com)
3. Access to the client ID, client secret, and an appropriately scoped grant token for the aforementioned application

### For Developers

If you're interested in contributing to the improvement, you'll need the tools to do so.

1. [Go v1.16](https://golang.org) - other versions may work, but I've only tested on v1.16

    There are other ways, of course, but for maximum flexibility, I installed Go using [asdf](https://asdf-vm.com). See the installation instructions below.

## Installation

1. Download a released binary and place it into your `$PATH`
1. Run `site24x7 configure` to provide authentication and authorization credentials
1. `site24x7 --help` to see what's available

## Development

1. Clone this repository
1. Create a `.env` file and supply appropriate values

        cp .env.example .env

1. Install project dependencies

        go install

1. Create a configuration file, either manually (see the _Configuration_ section below), or by building the tool and running `./site24x7 configure`.
1. Write code, submit pull requests!

## Configuration

The tool uses an external configuration file for authentication and authorization values and environment variables for less secret, but commonly used values.
