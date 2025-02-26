package md

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const input = `
---
title: goldmark-frontmatter
tags: [markdown, goldmark, working]
description: |
  Adds support for parsing YAML front matter.
---
# Hello

This is a test.

$$
\mathbb{E}(X) = \int x d F(x) = \left\{ \begin{aligned} \sum_x x f(x) \; & \text{ if } X \text{ is discrete} 
\\ \int x f(x) dx \; & \text{ if } X \text{ is continuous }
\end{aligned} \right.
$$


Inline math $\frac{1}{2}$
`

func TestParse(t *testing.T) {
	actual, err := Parse("hello.md", []byte(input))
	t.Log(actual)
	assert.NoError(t, err)
	assert.NotEmpty(t, actual)
	// t.Fail()
}
