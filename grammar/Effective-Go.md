## Formatting

The go provides the commands program of gofmt to formatting the go code style.
which operates at the package level rather than source file level.

```shell
go fmt
```

## Commentary

Go provides C-style /\* \*/ block comments and C++-style // line comments.

- Every package should have a package comment,a block comment preceding the package clause.
- For multi-file packages, the package comment only needs to be present in one file, and any one will do.
- If the package is simple, the package comment can be brief.
- Every exported (capitalized) name in a program should have a doc comment.
- The first sentence should be a one-sentence summary that starts with the name being declared.

If every doc comment begins with the name of the item it describes, you can use the doc subcommand of the go tool and run the output through grep.

```shell
go doc -all regexp | grep -i item-name
```

## Package names

Notes: Don't use the import . notation, which can simplify tests that must run outside the package they are testing, but should otherwise be avoided.

## Getters
