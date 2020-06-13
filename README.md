# Monkey2

`Monkey2` is a program to generate code automatically, and run a sidecar to watch these code.

`insane` monkey generates code randomly.

`cheating` monkey copies code from the repo given thru env **CHEATING_REPO**. 

## Install & Run

```shell script
❯ go install github.com/git-roll/monkey2
❯ monkey -h
monkey [name] [sidecar]

name could be one of [insane, cheating].
You can also run a sidecar to watch the monkey. e.g.

> monkey insane git roll

> monkey cheating git roll
```
