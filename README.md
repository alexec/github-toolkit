# Github Toolkit

Install:

```
GO111MODULE=on go install github.com/alexec/github-toolkit/cmd/gt
```

Create release note:

```
cd my-repo
gt relnote v1.3.0-rc3..v1.3.0-rc4
```

Create cards:

```
cd my-repo
./gt cards --help
```

![cards](docs/images/cards.png)


# Building

```
make
```