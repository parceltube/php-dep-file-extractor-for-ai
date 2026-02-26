<h1 align="center">PHP Dependency Extractor</h1>

<p align="center">
  Extract selected PHP files and their dependencies from large projects in minutes.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Platform-Windows-0078D6?style=flat-square" alt="Windows">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Frameworks-ZF1%20%7C%20CakePHP%20%7C%20Laravel-6f42c1?style=flat-square" alt="Frameworks">
  <img src="https://img.shields.io/badge/License-See%20LICENSE.txt-brightgreen?style=flat-square" alt="License">
</p>

<p align="center">
  <a href="./manual.html"><b>English Manual</b></a> •
  <a href="./manual-cn.html"><b>中文手册</b></a>
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> |
  <a href="#user-documentation">Docs</a> |
  <a href="#install-and-build">Build</a> |
  <a href="#project-structure">Structure</a> |
  <a href="#release-for-end-users">Release</a> |
  <a href="#notes">Notes</a> |
  <a href="#license">License</a>
</p>

---

## Why This Tool

`PHP Dependency Extractor` helps you isolate a minimal and relevant code bundle from legacy or large PHP codebases.

Typical use cases:
- Share only related files with teammates for code review
- Prepare a focused code package for AI-assisted analysis
- Extract one controller/service and its dependencies for debugging

## Features

- Scan all `.php` files and build a searchable file tree
- Resolve class dependencies from selected files
- Framework-aware indexing and resolution:
  - ZF1
  - CakePHP
  - Laravel
- Optional `require/include` parsing with manual selection
- Preserve original relative folder structure on export
- Local-only server binding (`127.0.0.1`)
- Single executable runtime experience (Windows)

## Quick Start

1. Run `php-dep-extractor.exe`
2. Browser opens automatically at `http://127.0.0.1:<port>`
3. Choose project folder and framework (`ZF1` / `CakePHP` / `Laravel`)
4. Click `Scan`
5. Select files in the tree
6. Click `Analyze`
7. Click `Copy Files`

## User Documentation

English:
- [manual.html](./manual.html)
- [manual.md](./manual.md)

Chinese:
- [manual-cn.html](./manual-cn.html)
- [manual-cn.md](./manual-cn.md)

## Install and Build

Build from source:

```bash
go build -o php-dep-extractor.exe .
```

Run:

```bash
./php-dep-extractor.exe
```

## Project Structure

```text
.
├─ main.go
├─ internal/
│  ├─ server/     # HTTP handlers and app state
│  ├─ scanner/    # Project scan and class index
│  ├─ parser/     # Dependency/include parsing
│  ├─ filetree/   # Tree builder
│  └─ copier/     # File export logic
├─ web/
│  ├─ index.html
│  ├─ app.js
│  └─ style.css
├─ manual.html
├─ manual-cn.html
├─ manual.md
└─ manual-cn.md
```

## Release for End Users

If your users should "download and run directly", publish binaries in **GitHub Releases**:

- Upload `php-dep-extractor.exe` (or a zip package)
- Add links to `manual.html` / `manual-cn.html` in release notes
- Keep source repository clean (do not commit `.exe` into source tree)

## Notes

- Current UX is primarily designed for Windows (folder picker uses PowerShell)
- Dependency detection is regex/rule-based, not full AST semantic parsing

## License

This project is licensed under the terms in [LICENSE.txt](./LICENSE.txt).
