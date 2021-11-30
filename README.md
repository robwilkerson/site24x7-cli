# Site24x7 Command Line Client

A small CLI tool for interacting with the [Site24x7 API](https://www.site24x7.com/help/api/#introduction).


This is currently _unreleased software_. The content of this README is anticipatory more than actual.

## Prerequisites

### For Basic Use

Theoritically none, but I haven't yet tested anything against architectures _other_ than the one I'm running.

### For Developers

If you're interested in contributing to the improvement, you'll need the tools to do so.

1. [Go v1.16](https://golang.org)

    There are other ways, of course, but for maximum flexibility, I installed Go using [asdf](https://asdf-vm.com). See the installation instructions below.

## Installation

### For Basic Use

1. Download the latest `site24x7` binary
2. Save the binary to a directory in your `$PATH`
3. Run `site24x7 --help` to see what there is to see

### For Developers

As mentioned above, these instructions will assume you're installing Go using asdf. If you choose another path, that's cool too.

### Using asdf

1. Clone this repository
1. Create a `.env` file and supply redacted values

        cp .env.example .env

1. Install the Go plugin for asdf

        asdf plugin add golang

1. Install Go

        asdf install golang 1.16.7

1. Make sure that Go is in your `$PATH`

        echo PATH=$PATH:$(go env GOPATH)/bin >> <your shell startup file>

1. Set the project's default version of Go

        cd <project directory> && asdf local golang 1.16.7

1. Verify the Go version in your project directory

        go version

1. Install project dependencies

        go install

1. Write code, submit pull requests!
