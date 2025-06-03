---
id: lzma-go
aliases: []
tags:
  - programming-language/go
  - ideology/compression/lzma
banner_path: /dist/img/projects/lzma-go.webp
created_at: 2025-01-06T12:00:00.000-06:00
description: Pure Go implementation of LZMA compression algorithm
slug: lzma-go
title: lzma-go
updated_at: 2025-06-03T14:58:32.000-06:00
---

A pure Go implementation of the LZMA (Lempel-Ziv-Markov chain Algorithm) compression algorithm. This library provides efficient compression and decompression capabilities without requiring CGO or external dependencies.

## Features

- Pure Go implementation - no CGO required
- Compatible with standard LZMA format
- Efficient compression ratios
- Streaming compression/decompression support
- Thread-safe operations

## Installation

```bash
go get github.com/conneroisu/lzma-go
```

## Usage

```go
import "github.com/conneroisu/lzma-go"
```

Package lzma package implements reading and writing of LZMA format compressed data.

Reference implementation is LZMA SDK version 4.65 originally developed by Igor Pavlov, available online at:

```
http://www.7-zip.org/sdk.html
```

Usage examples. Write compressed data to a buffer:

```
var b bytes.Buffer
w := lzma.NewWriter(&b)
w.Write([]byte("hello, world\n"))
w.Close()
```

read that data back:

```
r := lzma.NewReader(&b)
io.Copy(os.Stdout, r)
r.Close()
```

If the data is bigger than you'd like to hold into memory, use pipes. Write compressed data to an io.PipeWriter:

```
pr, pw := io.Pipe()
 go func() {
 	defer pw.Close()
	w := lzma.NewWriter(pw)
	defer w.Close()
	// the bytes.Buffer would be an io.Reader used to read uncompressed data from
	io.Copy(w, bytes.NewBuffer([]byte("hello, world\n")))
 }()
```

and read it back:

```
defer pr.Close()
r := lzma.NewReader(pr)
defer r.Close()
// the os.Stdout would be an io.Writer used to write uncompressed data to
io.Copy(os.Stdout, r)
```

## Performance

The library achieves competitive compression ratios while maintaining reasonable performance for a pure Go implementation. It's suitable for applications where CGO dependencies need to be avoided.

## Links

- [GitHub Repository](https://github.com/conneroisu/lzma-go)
