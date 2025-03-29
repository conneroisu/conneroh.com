---
id: seltabl
aliases: []
tags:
  - programming-language/go
  - programming-language/html
created_at: 2025-03-28T19:12:39.000-06:00
description: A Go library for extracting data from HTML tables.
title: seltabl
updated_at: 2025-03-28T19:19:02.000-06:00
---

Embarking on a journey through the intricate world of web data extraction, I often found myself entangled in the complexities of parsing HTML tables. The repetitive tasks, the cumbersome code, and the constant battle to maintain efficiency were challenges I knew many developers faced. Determined to find a solution, I channeled my passion for coding and problem-solving into creating **seltabl**, a Go library designed to transform the way we interact with HTML tables.

## The Genesis of seltabl

The inception of seltabl was driven by a simple yet profound realization: extracting data from HTML tables shouldn't be a herculean task. As a senior in Electrical Engineering and Computer Science at Iowa State University, I had encountered numerous scenarios where the need for a streamlined, efficient method to parse HTML tables was evident. Leveraging the power of Go and inspired by the capabilities of the goquery library, I set out to develop a tool that would not only simplify the process but also offer configurability and robust developer support.

## Unveiling seltabl

At its core, seltabl is a Go library accompanied by a command-line interface (CLI) and a language server. Its primary function is to parse HTML sequences into structs, making it particularly adept at handling HTML tables. However, its versatility allows it to be employed for any HTML sequence parsing tasks. By enabling data binding to structs and providing a dynamic way to define table schemas, seltabl empowers developers to interact with web data more intuitively and efficiently.

## Key Features

- **Configurable Parsing:** seltabl allows for customizable parsing configurations, enabling developers to tailor the extraction process to their specific needs.
  
- **Goquery Integration:** By leveraging goquery, seltabl offers a familiar and powerful selection syntax, akin to jQuery, making it accessible for those already acquainted with CSS selectors.
  
- **Developer Tooling:** The inclusion of a language server and CLI utilities enhances the development experience, providing tools for code generation, linting, and testing.

## Getting Started with seltabl

Integrating seltabl into your Go project is straightforward. Begin by installing the package:

```bash
go get github.com/conneroisu/seltabl
```

For access to the CLI and language server functionalities, install the command-line tool:

```bash
go install github.com/conneroisu/seltabl/tools/seltabls@latest
```

## A Glimpse into Usage

Imagine you have an HTML table and you wish to extract its data into Go structs. With seltabl, this task becomes seamless. Here's a basic example:

```go
package main

import (
    "fmt"
    "github.com/conneroisu/seltabl"
)

type TableRow struct {
    Column1 string `json:"column1" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" ctl:"text"`
    Column2 string `json:"column2" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" ctl:"text"`
}

func main() {
    htmlContent := `
    <table>
        <tr>
            <td>Data1</td>
            <td>Data2</td>
        </tr>
        <tr>
            <td>Data3</td>
            <td>Data4</td>
        </tr>
    </table>
    `
    rows, err := seltabl.NewFromString[TableRow](htmlContent)
    if err != nil {
        panic(fmt.Errorf("failed to parse HTML: %w", err))
    }
    for _, row := range rows {
        fmt.Printf("%+v\n", row)
    }
}
```

This snippet demonstrates how seltabl can be employed to parse an HTML table, extracting its contents into a slice of `TableRow` structs. The use of struct tags like `hSel`, `dSel`, and `ctl` allows for precise selection and extraction of data, showcasing the library's configurability.

## The Road Ahead

The development of seltabl has been a journey of learning, innovation, and community engagement. As I continue to refine and expand its capabilities, I invite fellow developers to explore, contribute, and provide feedback. Together, we can make web data extraction in Go not just a necessity, but a delight.

For more information, detailed documentation, and contribution guidelines, visit the [seltabl GitHub repository](https://github.com/conneroisu/seltabl). Let's embark on this journey together, transforming challenges into opportunities and code into solutions. 
