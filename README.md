# Github Toolkit

[![Build Status](https://travis-ci.org/alexec/github-toolkit.svg?branch=master)](https://travis-ci.org/alexec/github-toolkit)

Install:

```
GO111MODULE=on go install github.com/alexec/github-toolkit/cmd/ght
```

## Release Notes

Creates a release note but examining the commits as follows:
 
1. If there is an issues ID in the commit message, that is used.
1. If there is nothing in the message, the the commit goes in an "other" bucket.

If the issues is actually a pull request, we check the pull request body for normal issues. 

For each issue ID that is found, they are categorised as follows:
 
* Enhancement - if the issue is labelled with "enhancement".
* Bug fix - if the issue is labelled with "bug".
* Other - otherwise

Create release note:

```
ght relnote v1.3.0-rc3..v1.3.0-rc4
```


## Issue Cards

Github does not provide a way to generate arbitrary issue cards for agile planning. This command creates a HTML page which list issues and can be printed: 

```
# enhancements backlog 
ght cards --label enhancement --exclude-label wontfix --milestone none 

# bugs backlog
ght cards --label bug --exclude-label wontfix --milestone none 

# help wanted backlog
ght cards --label 'help wanted' --exclude-label wontfix' --milestone none 

# open issues in milestone v1.3
ght cards --milestone v1.3

# issues opened in the last day
ght cards --state all --since 24h
```

![cards](docs/images/cards.png)


# Building

```
cd ~/go/src/github.com/alexec/github-toolkit
make
```
