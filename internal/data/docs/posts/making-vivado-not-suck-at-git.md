---
id: making-vivado-not-suck-at-git
aliases:
  - Making Vivado not suck at Git
tags:
  - ide/vivado
  - vcs/git
  - edu/iastate/cpre488
created_at: 2025-03-27T14:13:10.000-06:00
description: Post on how I made Vivado not suck at Git
projects:
  - cpre488-mp2
title: Making Vivado not suck at Git
updated_at: 2025-03-28T20:07:29.000-06:00
---

# Making Vivado not suck at Git

Vivado is a great tool for prototyping and simulation, but it's not really
great for version control. I've been using Vivado for a while now inside of 
[CPRE 488](/tags/edu/iastate/cpre488), and I've found that it's really
difficult to use Git with Vivado.

## The Problem

Vivado creates a lot of files 
