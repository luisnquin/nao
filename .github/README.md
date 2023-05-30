
<h1 align="center"><code>nao🍵</code></h1>

<div align="center">

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/luisnquin/nao)](https://github.com/luisnquin/nao)
[![CI](https://github.com/luisnquin/nao/actions/workflows/go.yml/badge.svg)](https://github.com/luisnquin/nao/actions/workflows/go.yml)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![GitHub stars](https://img.shields.io/github/stars/luisnquin/nao.svg?style=social&label=Star&maxAge=2592000)](https://github.com/luisnquin/nao)
[![Built with Nix](https://img.shields.io/static/v1?logo=nixos&logoColor=white&label=&message=Built%20with%20Nix&color=41439a)](https://github.com/luisnquin/nao)

<p>Take notes without worrying about the path where the file is</p>
</div>

## Features

- [x] You know terminal, you know nao
- [x] No need to specify a path to access a note
- [x] Edit from terminal editor
- [x] One writer and multiple readers by note

## Demo

[![asciicast](https://asciinema.org/a/9DETM5MtJaA9d0emviPvz1n0s.svg)](https://asciinema.org/a/9DETM5MtJaA9d0emviPvz1n0s)

## Install

```bash
# Requires go 1.18>=
$ go install github.com/luisnquin/nao/v3/cmd/nao@v3.1.4
```

## Completions

Add the line(s) of your corresponding shell to your .zshrc|.bashrc file

```bash
# bash
source <(nao completion bash)

# zsh
source <(nao completion zsh)
compdef _nao nao
```

## Why did I do this?

No one has been able to do this, so here we are

## License

[MIT](https://raw.githubusercontent.com/luisnquin/nao/main/LICENSE)
