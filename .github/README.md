# Nao - notes manager, writing has become too comfortable in the 21st century

[![Go](https://github.com/luisnquin/nao/actions/workflows/go.yml/badge.svg)](https://github.com/luisnquin/nao/actions/workflows/go.yml)

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

## Installation

- Via Go install

    ```bash
        # Requires go 1.18>=
        $ go install github.com/luisnquin/nao/v3/cmd/nao@latest
    ```

- Via [Nix](https://nix.dev/)

    ```bash
        nix-env -iA nixpkgs.nao
    ```

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

In construction

## License

MIT © [Luis Quiñones](https://github.com/luisnquin)
