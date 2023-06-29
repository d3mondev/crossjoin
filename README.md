# Crossjoin

![Example output](https://i.imgur.com/XeJQCqp.png)

## Description

`crossjoin` is a utility designed for security testing and fuzzing tasks. It takes multiple input files and creates all possible permutations (the Cartesian product) of their lines. This can be useful for generating comprehensive lists of potential URLs or other input data for fuzzing and penetration testing purposes.

For example, with input files containing HTTP protocols, domain names, and URL paths respectively, `crossjoin` will generate all possible combinations of these components, thereby creating a comprehensive list of URLs.

[![Twitter URL](https://img.shields.io/twitter/url/https/twitter.com/d3mondev.svg?style=social&label=Follow%20%40d3mondev)](https://twitter.com/d3mondev)

## Installation
You can download a release from the [Releases](https://github.com/d3mondev/crossjoin/releases) page.

Alternatively, you can compile it yourself. You need to have [Go](https://go.dev/dl/) installed on your system.

```
go install github.com/d3mondev/crossjoin@latest
```

## Usage

```bash
crossjoin file1 file2 file3 [...fileN]
```

Each file should contain a set of strings (lines) to be used in the permutations. `crossjoin` will then output the permutations to the console.

If standard input (stdin) is provided, the program will use it as the first input.

```bash
command | crossjoin file1 file2 ...
```

## Example

Given the following 3 files:

#### protocols.txt:
```plaintext
http://
https://
```

#### domains.txt:
```plaintext
example.com
www.example.com
```

#### paths.txt:
```plaintext
/index.html
/admin/
```

Running `crossjoin` with these files as input will produce:

```
$ crossjoin protocols.txt domains.txt paths.txt

http://example.com/index.html
http://example.com/admin/
http://www.example.com/index.html
http://www.example.com/admin/
https://example.com/index.html
https://example.com/admin/
https://www.example.com/index.html
https://www.example.com/admin/
```
