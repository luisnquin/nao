
![anyway-removebg-preview (1)](https://user-images.githubusercontent.com/86449787/216045021-2a8d3a1c-1bbc-497b-97aa-74c9dc399616.png)

# nao

[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/luisnquin/nao)](https://github.com/luisnquin/nao)
[![CI](https://github.com/luisnquin/nao/actions/workflows/go.yml/badge.svg)](https://github.com/luisnquin/nao/actions/workflows/go.yml)
[![GitHub stars](https://img.shields.io/github/stars/luisnquin/nao.svg?style=social&label=Star&maxAge=2592000)](https://github.com/luisnquin/nao)
![Love](https://img.shields.io/badge/Love-pink?style=flat-square&logo=data:image/svg%2bxml;base64,PHN2ZyByb2xlPSJpbWciIHZpZXdCb3g9IjAgMCAyNCAyNCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48dGl0bGU+R2l0SHViIFNwb25zb3JzIGljb248L3RpdGxlPjxwYXRoIGQ9Ik0xNy42MjUgMS40OTljLTIuMzIgMC00LjM1NCAxLjIwMy01LjYyNSAzLjAzLTEuMjcxLTEuODI3LTMuMzA1LTMuMDMtNS42MjUtMy4wM0MzLjEyOSAxLjQ5OSAwIDQuMjUzIDAgOC4yNDljMCA0LjI3NSAzLjA2OCA3Ljg0NyA1LjgyOCAxMC4yMjdhMzMuMTQgMzMuMTQgMCAwIDAgNS42MTYgMy44NzZsLjAyOC4wMTcuMDA4LjAwMy0uMDAxLjAwM2MuMTYzLjA4NS4zNDIuMTI2LjUyMS4xMjUuMTc5LjAwMS4zNTgtLjA0MS41MjEtLjEyNWwtLjAwMS0uMDAzLjAwOC0uMDAzLjAyOC0uMDE3YTMzLjE0IDMzLjE0IDAgMCAwIDUuNjE2LTMuODc2QzIwLjkzMiAxNi4wOTYgMjQgMTIuNTI0IDI0IDguMjQ5YzAtMy45OTYtMy4xMjktNi43NS02LjM3NS02Ljc1em0tLjkxOSAxNS4yNzVhMzAuNzY2IDMwLjc2NiAwIDAgMS00LjcwMyAzLjMxNmwtLjAwNC0uMDAyLS4wMDQuMDAyYTMwLjk1NSAzMC45NTUgMCAwIDEtNC43MDMtMy4zMTZjLTIuNjc3LTIuMzA3LTUuMDQ3LTUuMjk4LTUuMDQ3LTguNTIzIDAtMi43NTQgMi4xMjEtNC41IDQuMTI1LTQuNSAyLjA2IDAgMy45MTQgMS40NzkgNC41NDQgMy42ODQuMTQzLjQ5NS41OTYuNzk3IDEuMDg2Ljc5Ni40OS4wMDEuOTQzLS4zMDIgMS4wODUtLjc5Ni42My0yLjIwNSAyLjQ4NC0zLjY4NCA0LjU0NC0zLjY4NCAyLjAwNCAwIDQuMTI1IDEuNzQ2IDQuMTI1IDQuNSAwIDMuMjI1LTIuMzcgNi4yMTYtNS4wNDggOC41MjN6Ii8+PC9zdmc+)
[![Built with Nix](https://img.shields.io/static/v1?logo=nixos&logoColor=white&label=&message=Built%20with%20Nix&color=41439a)](https://github.com/luisnquin/nao)

## Installation

- Via Go install

    ```bash
        # Requires go 1.18>=
        $ go install github.com/luisnquin/nao/v3/cmd/nao@v3.0.0
    ```

- From source

    ```bash
        # Requires git and go 1.18>=
        $ git clone https://github.com/luisnquin/nao.git
        $ cd nao
        $ make build
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

```bash
 # Type
 nao config
 # or to see more options that you can use to edit the configuration file
 nao config --help
```

## License

MIT © [Luis Quiñones](https://github.com/luisnquin)
