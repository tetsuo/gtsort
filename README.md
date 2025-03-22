# tsort-go

A Go port of the [`tsort`](https://github.com/coreutils/coreutils/blob/cb2abbac7f9e40e0f0d6183bf9b11e80b0cad8ef/src/tsort.c) utility from GNU Coreutils.
It performs a topological sort using Algorithm T from Knuth's *The Art of Computer Programming*, Volume 1.

## Usage

```bash
tsort-go [OPTION] [FILE]
```

Write totally ordered list consistent with the partial ordering in FILE. With no FILE, or when FILE is `-`, read standard input.

## Example

```
tsort-go <<EOF
a b c
d
e f
b c d e
EOF
```

will produce the output

```
a
b
c
d
e
f
```

For detailed information on the `tsort` command invocation, see the [GNU Coreutils manual](https://www.gnu.org/software/coreutils/manual/html_node/tsort-invocation.html).
