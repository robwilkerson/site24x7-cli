# Site24x7 Command Line Client

[![Build](https://github.com/robwilkerson/site24x7-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/robwilkerson/site24x7-cli/actions/workflows/ci.yml)

A CLI tool for interacting with the [Site24x7 API](https://www.site24x7.com/help/api/#introduction).

This is currently _early alpha-level software at best_. The content of this README is anticipatory more than actual.

## Prerequisites

1. A Site24x7 account
2. A registered [Site24x7 application](https://api-console.zoho.com)
3. Access to the client ID, client secret, and an appropriately scoped grant token for the aforementioned application

To register an application and generate a client ID, client secret, and grant token, follow the [API authentication instructions](https://www.site24x7.com/help/api/#authentication)

### For Developers

If you're interested in contributing to the improvement, you'll need the tools to do so.

1. [Go v1.17.5](https://golang.org) - other versions may work, but I've only tested on this one

    There are other ways, of course, but for maximum flexibility, I installed Go using [asdf](https://asdf-vm.com). See the installation instructions below.

## Installation

1. Download a released binary and place it into your `$PATH`
1. Run `site24x7 configure` to provide authentication and authorization credentials; you'll need to have your client ID, client secret, and grant token handy
1. `site24x7 --help` to see what's available

## Development

1. Clone this repository
1. Create a `.env` file and supply appropriate values

        cp .env.example .env

1. Install project dependencies

        go install

1. Create a configuration file by building the tool (`go build`) and running `./site24x7 configure`; you'll need to have your client ID, client secret, and grant token handy
1. Write code, test it, and submit pull requests!
