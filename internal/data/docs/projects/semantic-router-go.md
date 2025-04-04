---
id: semantic-router-go
aliases:
  - semanticrouter-go
tags:
  - programming-language/go
created_at: 2025-03-28T19:17:07.000-06:00
description: A high-performance, cost-effective AI decision-making library written in pure Go.
title: semantic-router-go
updated_at: 2025-03-28T20:07:30.000-06:00
---

# semanticrouter-go

As a college student with a keen interest in artificial intelligence and software development, I developed **semanticrouter-go**, a high-performance, cost-effective AI decision-making library written in pure Go. This project aims to enhance the efficiency of large language models (LLMs) and AI agents by providing a rapid decision-making layer that leverages semantic vector spaces for routing requests based on configurable semantic meanings. citeturn0search0

**Key Features of semanticrouter-go:**

- **Efficient Decision-Making:** By utilizing semantic vector spaces, the library enables swift routing of requests, eliminating the latency associated with traditional LLM-generated decisions.

- **Pure Go Implementation:** Written entirely in Go, semanticrouter-go ensures seamless integration into Go-based projects without external dependencies.

- **Flexible Encoding Support:** The library supports various encoding methods, including integration with models like `mxbai-embed-large` through the Ollama API, facilitating versatile embedding strategies.

- **Customizable Routing:** Developers can define routes with specific utterances, allowing the router to match incoming queries to the most appropriate route based on semantic similarity.

**Installation:**

To incorporate semanticrouter-go into your Go project, execute:

```bash
go get github.com/conneroisu/semanticrouter-go
```

**Example Use Case:**

Consider a scenario in a veterinary application where it's essential to distinguish between noteworthy medical inquiries and casual chitchat. semanticrouter-go can be configured to route user inputs accordingly:

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/semanticrouter-go"
	"github.com/conneroisu/semanticrouter-go/encoders/ollama"
	"github.com/conneroisu/semanticrouter-go/stores/memory"
	"github.com/ollama/ollama/api"
)

var NoteworthyRoutes = semanticrouter.Route{
	Name: "noteworthy",
	Utterances: []semanticrouter.Utterance{
		{Utterance: "What is the best way to treat a dog with a cold?"},
		{Utterance: "My cat has been limping; what should I do?"},
	},
}

var ChitchatRoutes = semanticrouter.Route{
	Name: "chitchat",
	Utterances: []semanticrouter.Utterance{
		{Utterance: "What's your favorite color?"},
		{Utterance: "Do you like animals?"},
	},
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	cli, err := api.ClientFromEnvironment()
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	router, err := semanticrouter.NewRouter(
		[]semanticrouter.Route{NoteworthyRoutes, ChitchatRoutes},
		&ollama.Encoder{
			Client: cli,
			Model:  "mxbai-embed-large",
		},
		memory.NewStore(),
	)
	if err != nil {
		return fmt.Errorf("error creating router: %w", err)
	}
	finding, p, err := router.Match(ctx, "How's the weather today?")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Found:", finding.Name)
	fmt.Println("Probability:", p)
	return nil
}
```

In this example, the router evaluates the input "How's the weather today?" and classifies it under the "chitchat" route, demonstrating its ability to discern between different types of user interactions. citeturn0search0

**Development and Contribution:**

semanticrouter-go is designed with a focus on performance and simplicity, aiming to provide developers with a tool to enhance AI decision-making processes efficiently. The project is open-source and licensed under the MIT License, inviting community contributions and collaboration.

For more information, to explore the source code, or to contribute to the project, visit the GitHub repository: [https://github.com/conneroisu/semanticrouter-go](https://github.com/conneroisu/semanticrouter-go)

By developing semanticrouter-go, I aim to empower Go developers to implement rapid and intelligent decision-making capabilities within their AI applications, fostering innovation and efficiency in the field of artificial intelligence.
