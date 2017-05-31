# Babelfish Tools

[![Build Status](https://travis-ci.org/bblfsh/tools.svg?branch=master)](https://travis-ci.org/bblfsh/tools)
[![codecov](https://codecov.io/gh/bblfsh/tools/branch/master/graph/badge.svg)](https://codecov.io/gh/bblfsh/tools)

Language analysis tools on top of Babelfish

## Build

### With docker

`make build`

### Without docker

`make build-internal`

## Usage

Babelfish Tools provides a set of tools built on top of Babelfish, to
see which tools are supported, run:

`bblfsh-tools --help`

To make use of any of these tools you need to have the Babelfish
server up and running. Look at
[server site](https://github.com/bblfsh/server/) for details.

Once you have a server running, you can use the dummy tool, which
should let you know if the connection with the server succeeded:

`bblfsh-tools dummy path/to/source/code`

If the server is in a different location, use the `address` parameter:

`bblfsh-tools dummy --address location:port path/to/source/code`

Once connection with the server is working fine, you can use any other
available tool in a similar way.

## License

GPLv3, see [LICENSE](LICENSE)
