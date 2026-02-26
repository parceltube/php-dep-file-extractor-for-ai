# PHP Dependency Extractor

PHP Dependency Extractor is a local desktop tool (Go + Web UI) for extracting selected PHP files and their dependencies from large projects, then exporting them into an isolated folder for review, debugging, or AI-assisted analysis.

## Use Cases

- Extract a minimal set of related files from a legacy PHP codebase
- Prepare a reproducible package for teammate/vendor code review
- Build a focused context bundle for AI tools

## Key Features

- Scans `.php` files and builds a selectable file tree
- Resolves dependencies from selected files
- Framework support: `ZF1`, `CakePHP`, `Laravel`
- Optional `require/include` parsing with manual selection
- Preserves original relative folder structure on export
- Local server binds to `127.0.0.1` only

## Quick Start

1. Run `php-dep-extractor.exe`
2. The browser opens automatically at `http://127.0.0.1:<port>`
3. Choose a project directory and framework, then click `Scan`
4. Select files in the left tree, then click `Analyze`
5. Click `Copy Files` to export

## Detailed User Guide

For full operation details, see:

- [manual.html](./manual.html)

If you are viewing on GitHub and prefer Markdown:

- [manual.md](./manual.md)

Chinese documentation:

- [manual-cn.html](./manual-cn.html)
- [manual-cn.md](./manual-cn.md)

## Build From Source

```bash
go build -o php-dep-extractor.exe .
```

## Project Structure

```text
.
├─ main.go
├─ internal/
│  ├─ server/
│  ├─ scanner/
│  ├─ parser/
│  ├─ filetree/
│  └─ copier/
├─ web/
│  ├─ index.html
│  ├─ app.js
│  └─ style.css
├─ manual.html
└─ manual.md
```

## Notes

- Current UX is primarily targeted for Windows (folder picker uses PowerShell)
- Dependency parsing is rule/regex-based, not full AST semantic analysis

## License

This project is licensed under the terms in [LICENSE.txt](./LICENSE.txt).
