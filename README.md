# pony

[![Travis CI](https://img.shields.io/travis/jessfraz/pony.svg?style=for-the-badge)](https://travis-ci.org/jessfraz/pony)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/pony)

Local file-based password, API key, secret, recovery code store backed by GPG.

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
NAME:
   pony - Local File-Based Password, API Key, Secret, Recovery Code Store Backed By GPG

USAGE:
   pony [global options] command [command options] [arguments...]

VERSION:
   version v0.2.1, build 33bfbcc

AUTHOR(S):
   @jessfraz <no-reply@butts.com>

COMMANDS:
   add, save    Add a new secret
   delete, rm   Delete a secret
   get          Get the value of a secret
   list, ls     List all secrets
   update       Update a secret
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d                  run in debug mode
   --file, -f "~/.pony"         file to use for saving encrypted secrets
   --gpgpath "~/.gnupg/"        filepath used for gpg keys
   --keyid                      optionally set specific gpg keyid/fingerprint to use for encryption & decryption [$PONY_KEYID]
   --help, -h                   show help
   --generate-bash-completion
   --version, -v                print the version
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
$ pony add com.twitter.frazelledazzell.token KJDHJKFHDSBJDF
# GPG Passphrase for key "Jess Frazelle <butts@systemd.lol>":

$ pony add com.github.jessfraz.token LKJHSDLFKJDHF
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
