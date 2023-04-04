
# nao 🍵

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/luisnquin/nao)](https://github.com/luisnquin/nao)
[![CI](https://github.com/luisnquin/nao/actions/workflows/go.yml/badge.svg)](https://github.com/luisnquin/nao/actions/workflows/go.yml)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![GitHub stars](https://img.shields.io/github/stars/luisnquin/nao.svg?style=social&label=Star&maxAge=2592000)](https://github.com/luisnquin/nao)
[![Built with Nix](https://img.shields.io/static/v1?logo=nixos&logoColor=white&label=&message=Built%20with%20Nix&color=41439a)](https://github.com/luisnquin/nao)

## Install

```bash
    # Requires go 1.18>=
    $ go install github.com/luisnquin/nao/v3/cmd/nao@v3.1.1
```

## Features

- [x] You know terminal, you know nao
- [x] No need to specify a path to access a note
- [x] Edit from terminal editor
- [x] Encryption
- [x] One writer by note, multiple readers

## Some basic commands

All the available commands will sound familiar to you

- **ls**

    ```bash
        # List all notes
        $ nao ls
    ```

    Output:

    ```console
        ID          TAG                LAST UPDATE    SIZE      VERSION
        2bcf4a4b0e  thoughts           6 minutes ago  356.01KB  611
        5ed4323ace  todo               5 hours ago    6.11KB    635
    ```

- **new**

    ```bash
        # automatically generates a random tag
        $ nao new                                                    INT ×
        
        # or

        # In this case, the tag will be 'thoughts'
        $ nao new thoughts
    ```

- **mod**

    ```bash
        # Opens your terminal editor, by default nano 🧤
        $ nao mod thoughts                                           INT × 

        # or

        $ nao mod 596ca14a32
    ```

- **cat**

    ```bash
        # Displays the note in stdout
        $ nao cat thoughts
    ```

    Output:

    ```console
        Seeing you in my tears
        In my own reflection
        I hear you in the wind that passes through me
        Feel you in my hunger
        You're haunting my ambition
        Beautifully destructive attraction
    ```

- **rm**

    ```bash
        # Deletes the note
        $ nao rm thoughts
    ```

    Output:

    ```console
    Are you sure you want to delete this note thoughts(2bcf4a4b0e/0.37KB)?: 
     ▸ Yes
       No

    # And then
    2bcf4a4b0e21483d85f2806cfbc03c8c

    ```

    Keep in mind, if you know docker then you know nao

## Configuration

```bash
 # Type
 nao config
 # or to see more options that you can use to edit the configuration file
 nao config --help
```

## License

MIT © [Luis Quiñones](https://github.com/luisnquin)
