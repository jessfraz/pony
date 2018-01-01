# pony

[![Travis CI](https://travis-ci.org/jessfraz/pony.svg?branch=master)](https://travis-ci.org/jessfraz/pony)

Local file-based password, API key, secret, recovery code store backed by GPG.

## Installation

#### Binaries

- **darwin** [386](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-darwin-386) / [amd64](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-darwin-amd64)
- **freebsd** [386](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-freebsd-386) / [amd64](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-freebsd-amd64)
- **linux** [386](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-linux-386) / [amd64](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-linux-amd64) / [arm](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-linux-arm) / [arm64](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-linux-arm64)
- **solaris** [amd64](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-solaris-amd64)
- **windows** [386](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-windows-386) / [amd64](https://github.com/jessfraz/pony/releases/download/v0.1.0/pony-windows-amd64)

#### Via Go

```bash
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
   version v0.1.0, build 33bfbcc

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
