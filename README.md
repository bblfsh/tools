# Babelfish Tools

Language analysis tools on top of Babelfish

## Usage

Babelfish Tools provides a set of tools built on top of Babelfish, to
see which tools are supported, run:

`bblfsh-tools --help`

There's a dummy tool which should let you know if the connection with
the server succeeded:

`bblfsh-tools dummy path/to/source/code`

If the server is in a different location, use the `address` parameter:

`bblfsh-tools dummy --address location:port path/to/source/code`

Once connection with the server is working fine, you can use any other
available tool in a similar way.
