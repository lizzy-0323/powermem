# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial implementation of PowerMem Go SDK
- Core Memory operations (Add, Search, Update, Delete, GetAll, DeleteAll)
- OceanBase vector storage support
- OpenAI LLM integration
- OpenAI Embedder integration
- Intelligent deduplication feature
- Ebbinghaus forgetting curve algorithm
- Multi-agent memory management
- Context support for all operations
- Goroutine-safe operations with RWMutex
- Environment variable configuration support
- JSON configuration file support
- Comprehensive examples (basic, advanced, multi-agent)
- Chinese and English documentation

### Features

- ✅ Core Memory API
- ✅ Vector search with OceanBase
- ✅ Intelligent deduplication
- ✅ Multi-agent support
- ✅ Metadata filtering
- ✅ Snowflake ID generation
- ✅ Ebbinghaus algorithm
- ✅ Configuration management
- ✅ Error handling
- ✅ Examples and documentation

### Upcoming

- [ ] SQLite storage support
- [ ] PostgreSQL + pgvector support
- [ ] Qwen LLM/Embedder support
- [ ] Anthropic LLM support
- [ ] Ollama support
- [ ] Hybrid search (vector + fulltext + sparse)
- [ ] Reranker support
- [ ] Graph storage support
- [ ] Multimodal support
- [ ] Unit tests
- [ ] Integration tests
- [ ] Benchmarks

## [0.1.0] - 2026-01-21

### Added

- Initial release of PowerMem Go SDK
- MVP implementation with core features

[Unreleased]: https://github.com/oceanbase/powermem-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/oceanbase/powermem-go/releases/tag/v0.1.0
