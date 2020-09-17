# pony

[![make-all](https://github.com/jessfraz/pony/workflows/make%20all/badge.svg)](https://github.com/jessfraz/pony/actions?query=workflow%3A%22make+all%22)
[![make-image](https://github.com/jessfraz/pony/workflows/make%20image/badge.svg)](https://github.com/jessfraz/pony/actions?query=workflow%3A%22make+image%22)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/pony)

Local file-based password, API key, secret, recovery code store backed by GPG.

**Table of Contents**

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Installation](#installation)
    - [Binaries](#binaries)
    - [Via Go](#via-go)
- [Usage](#usage)
  - [Best Practices](#best-practices)
    - [`HISTIGNORE`](#histignore)
    - [Namespacing Keys](#namespacing-keys)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/jessfraz/pony/releases).

#### Via Go

```console
$ go get github.com/jessfraz/pony
```

## Usage

```console
$ pony -h
pony -  Local File-Based Password, API Key, Secret, Recovery Code Store Backed By GPG.

Usage: pony <command>

Flags:

  -d, --debug  enable debug logging (default: false)
  --file       file to use for saving encrypted secrets (default: ~/.pony)
  --keyid      optionally set specific gpg keyid/fingerprint to use for encryption & decryption (or env var PONY_KEYID)

Commands:

  create   Create a secret.
  get      Get details for a secret.
  ls       List secrets.
  rm       Delete a secret.
  version  Show the version information.
```

### Best Practices

#### `HISTIGNORE`

You should obviously add pony to your `HISTIGNORE` for example:

```bash
export HISTIGNORE="ls:cd:cd -:pwd:exit:date:* --help:pony:pony *";
```

#### Namespacing Keys

You should namespace the keys for your secrets like the following:

```console
$ pony create com.twitter.frazelledazzell.token KJDHJKFHDSBJDF
# GPG Passphrase for key "Jess Frazelle <butts@systemd.lol>":

$ pony create com.github.jessfraz.token LKJHSDLFKJDHF
# GPG Passphrase for key "Jess Frazelle <butts@systemd.lol>":

# if a key ends with `.recovery`
# we assume it is a list of comma seperated
# strings that are recovery codes
$ pony add com.github.devnull@butts.com.recovery we0wk4,osdknew,4fd9kw,03jfn23,sduj39s
# GPG Passphrase for key "Jess Frazelle <butts@systemd.lol>":

$ pony ls
# GPG Passphrase for key "Jess Frazelle <butts@systemd.lol>":

KEY                                     VALUE
com.aws.amazon.prod.key                 KSUIIUEJDMSDBSDJFOFR
com.aws.amazon.prod.secret              skljdUYGjsndhfjjiosjdfgr/HKKSU
com.github.botaccount.recovery          we0wk4,osdknew,4fd9kw,03jfn23,sduj39s
com.github.jessfraz.token               LKJHSDLFKJDHF
com.twitter.frazelledazzell.token       KJDHJKFHDSBJDF

# you can also filter by a regular expression
$ pony ls --filter com.github*
# GPG Passphrase for key "Jess Frazelle <butts@systemd.lol>":

KEY                                     VALUE
com.github.botaccount.recovery          we0wk4,osdknew,4fd9kw,03jfn23,sduj39s
com.github.jessfraz.token               LKJHSDLFKJDHF
```
