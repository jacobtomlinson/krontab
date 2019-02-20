# krontab

[![Build Status](https://travis-ci.com/jacobtomlinson/krontab.svg?branch=master)](https://travis-ci.com/jacobtomlinson/krontab)
[![Current Release](https://img.shields.io/github/release/jacobtomlinson/krontab.svg)](https://github.com/jacobtomlinson/krontab/releases/latest)

A crontab replacement for kubernetes.

You can use it to create cron jobs on a kubernetes cluster in a familiar crontab format.
Krontab works by allowing you to create job templates which are used in kubernetes. Then create
specific cron jobs using the crontab. Example:

```
# template: default
0 1 * * * echo hello  # name: test
```

## Installation

### Linux (x86 64-bit)
```shell
curl -L https://github.com/jacobtomlinson/krontab/releases/download/{LATEST VERSION}/krontab-linux-amd64 -o /usr/local/bin/krontab
chmod +x /usr/local/bin/krontab
```

### Linux (arm 32-bit)
```shell
curl -L https://github.com/jacobtomlinson/krontab/releases/download/{LATEST VERSION}/krontab-linux-arm -o /usr/local/bin/krontab
chmod +x /usr/local/bin/krontab
```

### OS X (64-bit)
```shell
curl -L https://github.com/jacobtomlinson/krontab/releases/download/{LATEST VERSION}/krontab-darwin-amd64 -o /usr/local/bin/krontab
chmod +x /usr/local/bin/krontab
```

## Usage

```
$ krontab help
Krontab is a crontab replacement for kubernetes.

You can use it to create cron jobs on a kubernetes cluster in a familiar crontab format.
Krontab works by allowing you to create job templates which are used in kubernetes. Then create
specific cron jobs using the crontab. Example:

# Crontab example

# template: default
0 1 * * * echo hello  # name: test

Usage:
  krontab [flags]
  krontab [command]

Available Commands:
  create      Create a krontab resource
  delete      Delete a krontab resource
  edit        Edit a krontab resource
  help        Help about any command
  list        List krontab resources
  version     Print the version number of krontab

Flags:
  -e, --edit-crontab   Edit the crontab
  -h, --help           help for krontab
  -l, --list-crontab   List the crontab

Use "krontab [command] --help" for more information about a command.
```

## Contributing

### Environment

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
$ ./bin/krontab
```

#### Testing

``make test``
