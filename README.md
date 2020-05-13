# Monkey2

`Monkey2` is a program to generate code automatically, and run a sidecar to watch these code.

Only `insane` monkey is supported now, which generates code randomly.

## Install & Run

```shell script
❯ go install github.com/git-roll/monkey2
❯ monkey -h
monkey [name] [sidecar]

name could be one of [insane].
You can also run a sidecar to watch the monkey. e.g.

> monkey insane git roll
```
