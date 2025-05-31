---
id: bufnrix
aliases:
  - bufnrix
tags:
  - programming-language/go
  - programming-language/nix
  - programming-language/protobuf
  - protocol/grpc
banner_path: projects/bufrnix.webp
created_at: 2025-03-27T14:13:10.000-06:00
description: Nix powered protobuf tools
title: bufnrix
updated_at: 2025-05-31T11:53:02.000-06:00
---

# bufnrix

Nix powered Protocol Buffers with declarative, reproducible code generation and comprehensive developer tooling. 

## Why Bufrnix?

Protocol Buffer tooling has traditionally suffered from dependency hell, network dependencies, and non-reproducible builds. While Buf's remote plugin system simplifies initial setup, it introduces critical limitations that become deal-breakers for production teams:
### The Problems with Remote Plugin Systems

#### 🌐 Network Dependency Friction

- Remote plugins require constant internet connectivity, breaking offline development
- Corporate firewalls and air-gapped environments can't access remote plugin execution
- Network latency and rate limiting slow down development workflows
- Timeout errors (context deadline exceeded) and service interruptions disrupt CI/CD pipelines
- Geographic latency affects teams in regions distant from Buf's servers

#### 🔒 Security and Compliance Concerns

- Proprietary Protocol Buffer schemas must be sent to external servers for processing
- Financial services, healthcare, and government contractors can't share sensitive API definitions
- Intellectual property concerns prevent many organizations from using remote execution
- Compliance requirements (SOX, HIPAA, FedRAMP) demand local processing of technical specifications
- Supply chain security policies prohibit external dependency on third-party infrastructure

#### ⚡ Technical Limitations of Remote Plugin Systems

- 64KB response size limits cause silent failures with large generated outputs (affects protoc-gen-grpc-swift and other plugins)
- Plugins requiring file system access or multi-stage generation cannot function remotely
- "All" strategy requirement prevents efficient directory-based generation optimizations
- Custom plugins require expensive Pro/Enterprise subscriptions
- Plugin ecosystem growth is bottlenecked by centralized approval processes
- Cross-plugin dependencies (like protoc-gen-gotag modifying generated Go code) are impossible

#### 🔄 Reproducibility Challenges

- Network variability introduces non-determinism in generated code
- Plugin version updates can break existing workflows without warning
- Cache invalidation and remote infrastructure changes affect build consistency
- Migration between plugin versions often requires extensive code modifications
- Alpha-to-stable transitions have caused breaking changes requiring full codebase updates
- Remote caching can mask non-deterministic plugin behavior until production

### How Bufrnix Solves These Problems

#### 🏠 Local, Deterministic Execution

```nix
# All plugins execute locally with dependencies managed by Nix
languages.go = {
  enable = true;
  grpc.enable = true;     # No network calls, no timeouts
  validate.enable = true; # Full plugin ecosystem available
  # Exact plugin versions cryptographically pinned
  grpc.package = pkgs.protoc-gen-go-grpc; # v1.3.0 always
};
```

#### 🔐 Complete Privacy and Control

- All processing happens on your machines - schemas never leave your environment
- No external dependencies for code generation workflows
- Full control over plugin versions, updates, and security patches
- Compliance-friendly for regulated industries (SOX, HIPAA, FedRAMP)
- Supply chain integrity through cryptographic verification

#### ⚡ Performance and Flexibility

- 60x faster builds in some cases (20 minutes → 20 seconds in CI)
- No artificial size limits (64KB) or plugin capability restrictions
- Support for custom plugins, multi-stage generation, and complex workflows
- Plugin chaining and file system access work seamlessly
- Directory-based generation strategies for optimal performance
- Parallel execution across multiple languages and plugins

#### 🎯 True Reproducibility

```nix
# Same inputs = identical outputs, always
config = {
  languages.go.grpc.package = pkgs.protoc-gen-go-grpc; # Exact version pinned
  # Cryptographic hashes ensure supply chain integrity
  # Content-addressed storage prevents version drift
  # Hermetic builds with no external state
};
```

## Links

- [GitHub](https://github.com/conneroisu/bufnrix)

