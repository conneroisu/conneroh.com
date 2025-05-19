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
updated_at: 2025-05-16T09:33:21.000-06:00
---

# Making Vivado not suck at Git

Vivado is a great tool for prototyping and simulation, but it's not really
great for version control. I've been using Vivado for a while now inside of
[CPRE 488](/tags/edu/iastate/cpre488), and I've found that it's really
difficult to use Git with Vivado.

## The Problem

Vivado creates a lot of files

## Our Solution

You may need to modify the tcl script to include the correct files from a static location in the repo.

For example, in our case ([cpre488-mp2](/projects/cpre488-mp2)), we had to statically define places for:
- Constraints File
- VHDL Files (`design_1_wrapper.vhd`)
- HW IP Files (`avnet_hdmi_out`, `avnet_hdmi_in`, `interfaces`, `onsemi_vita_cam`, `onsemi_vita_spi`)

Tree view of the hardware pipelined vivado folder structure (with some files removed for clarity)
```bash
.
├── Constraints
│   └── master.xdc
├── design_1.pdf
├── digital_camera_pipeline.tcl
├── hw
│   └── IP
│       ├── avnet_hdmi_in
│       ├── avnet_hdmi_out
│       ├── interfaces
│       ├── onsemi_vita_cam
│       └── onsemi_vita_spi
└── VHDL
    └── design_1_wrapper.vhd
```



# Fixing Vivado Git

Add and Commit tcl script.

```bash
git add digital_camera_pipeline.tcl
git commit -m "added digital_camera_pipeline.tcl"
```

Upon new clone of repo, cd to repo directory and execute generate tcl script.

```bash
git clone https://github.com/conneroisu/cpre488-mp2.git
# IN VIVADO TCL CONSOLE -> 
# On IAState Windows PC
# cd C:/Users/connero/Downloads/cpre488-mp2/Vivado/digital_camera_pipeline

# IN VIVADO TCL CONSOLE -> 
# source digital_camera_pipeline.tcl
## OR
# TOOLS -> 
# Run Tcl Script (opens explorer then select tcl script)
```

Vivado Folder `.gitignore`
```gitignore
!*
.Xil/*
project_1/*
digital_camera/*
*.str
```
