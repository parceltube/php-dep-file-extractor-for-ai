# PHP Dependency Extractor - User Manual

## Overview

PHP Dependency Extractor (PDE) is a standalone Windows tool for extracting PHP source files and their dependencies from large projects. It automatically discovers class dependencies through framework naming conventions and regex parsing, then copies everything to an isolated folder for analysis.

**Use Case**: You have a 5000+ file PHP project and need to extract a small subset of related files (e.g., a controller and all its dependencies) to share with an AI assistant or for code review.

---

## Quick Start

1. Double-click `php-dep-extractor.exe`
2. Browser opens automatically at `http://127.0.0.1:<port>`
3. Click **Browse** next to Project to select your PHP project directory
4. Select a **Framework** (ZF1 / CakePHP / Laravel)
5. Click **Scan** to index the project
6. Check files in the tree you want to extract
7. Click **Analyze** to discover dependencies
8. Click **Copy Files** to export

---

## Interface

### Toolbar

| Control | Description |
|---------|-------------|
| **Project** | Path to your PHP project root. Click Browse to select. |
| **Output** | Where to copy extracted files. If left empty, auto-generates a sibling folder: `{project}_output_{timestamp}`. Double-click the field to reset to auto mode. |
| **Framework** | Select your framework for correct class name resolution: ZF1, CakePHP, or Laravel. |
| **Parse require/include** | Enable optional parsing of `require`/`include` statements. Results appear in a separate section for manual selection. |
| **Scan** | Traverse the project directory, build the file tree and class name index. |
| **Analyze** | Parse selected files for class references and resolve dependencies. |
| **Copy Files** | Copy all selected files + dependencies to the output directory. |
| **Settings** | Open settings panel (theme, font size, framework prefix mappings). |

### Left Panel — File Tree

- Displays all `.php` files found in the project
- **Collapsed by default** — click a folder to expand one level at a time
- Click a **checkbox** (or click the file row) to select/deselect files
- Use the **search box** to filter files by path — matching directories auto-expand

### Right Panel — Dependencies

After clicking Analyze, three sections appear:

| Section | Color | Description |
|---------|-------|-------------|
| **Selected** | Blue | Files you manually checked in the tree |
| **Dependencies** | Orange | Auto-discovered class dependencies, showing which class and reference type (new, extends, static, etc.) and which selected file references it |
| **Include/Require** | Gray | Only shown when "Parse require/include" is enabled. Each entry has a checkbox — check the ones you want to include in the copy |

### Status Bar

- Left: current operation status
- Right: file counts (Selected / Dependencies / Total) and version number

---

## Settings

Click the **Settings** button to open the settings panel. It has three tabs:

### General

- **Theme**: Choose between Dark, Light, or System (follows your Windows appearance setting)
- **Font Size**: Adjust from 11px to 18px using the slider or +/- buttons. Default is 13px. Setting persists across sessions.

### Framework

Configure **ZF1 prefix mappings** — the rules that map class name prefixes to directories under `application/`:

| Prefix | Directory | Example |
|--------|-----------|---------|
| `Parent_` | `parents/` | `Parent_DbTable_Car_X` → `parents/DbTable/Car/X.php` |
| `DbTable_` | `dbs/` | `DbTable_Car_CarrierCust` → `dbs/Car/CarrierCust.php` |
| `Service_` | `services/` | `Service_Sys_Breadcrumbs` → `services/Sys/Breadcrumbs.php` |
| `Model_` | `models/` | `Model_Car_CarrierCust` → `models/Car/CarrierCust.php` |
| `Form_` | `forms/` | `Form_Login` → `forms/Login.php` |

You can add, remove, or modify rows. Click **Save Mappings** to apply.

CakePHP and Laravel tabs show reference information about their detection methods.

### About

Version info and tool description.

---

## Framework Support

### ZF1 (Zend Framework 1)

**Class resolution**: Underscore-separated class names map to directory paths.

```
Model_Car_CarrierCust  →  application/models/Car/CarrierCust.php
DbTable_Ord_Order      →  application/dbs/Ord/Order.php
Service_Api_StarTrack  →  application/services/Api/StarTrack.php
```

**Detection patterns**:
- `new ClassName()`
- `extends ClassName`
- `implements InterfaceName`
- `ClassName::method()` (static calls)
- `function foo(ClassName $bar)` (type hints)

**Excluded**: Classes starting with `Zend_`, `ZendX_`, PHP built-ins (`stdClass`, `Exception`, `DateTime`, etc.)

### CakePHP

**Class resolution**: Based on CakePHP 2.x directory conventions.

| Type | Directory |
|------|-----------|
| Model | `Model/` |
| Controller | `Controller/` |
| Component | `Controller/Component/` |
| Behavior | `Model/Behavior/` |
| Helper | `View/Helper/` |

**Detection patterns**:
- `App::uses('ClassName', 'Type')`
- `App::import('Type', 'ClassName')`

### Laravel

**Class resolution**: PSR-4 namespace to path mapping.

```
App\Models\User             →  app/Models/User.php
App\Http\Controllers\OrderController  →  app/Http/Controllers/OrderController.php
App\Services\ShippingService          →  app/Services/ShippingService.php
```

**Detection patterns**:
- `use App\Models\User` statements
- Namespace-qualified class references

---

## Require/Include Parsing

When enabled via the checkbox, PDE scans selected files for:

- `require 'path'`
- `require_once 'path'`
- `include 'path'`
- `include_once 'path'`

**Supported path patterns**:

| Pattern | Example |
|---------|---------|
| Direct string | `require_once 'lib/helper.php'` |
| APPLICATION_PATH | `require_once APPLICATION_PATH . '/configs/constants.php'` |
| dirname(__FILE__) | `include dirname(__FILE__) . '/../bootstrap.php'` |
| __DIR__ | `require __DIR__ . '/functions.php'` |

Results appear in the gray "Include/Require" section. They are **not** automatically included — check the ones you want before copying.

---

## Output Structure

Copied files preserve their original directory structure relative to the project root:

```
Input project: C:\projects\myapp\
Selected:      application/controllers/V3/CustomersController.php
Dependencies:  application/services/Sys/Breadcrumbs.php
               application/dbs/Car/CarrierCust.php

Output: C:\projects\myapp_output_20260223_153045\
        └── application/
            ├── controllers/V3/CustomersController.php
            ├── services/Sys/Breadcrumbs.php
            └── dbs/Car/CarrierCust.php
```

---

## Fallback Class Detection

For files that don't match any framework naming convention (e.g., standalone utility classes, controllers), PDE reads the first 100 lines of the file looking for:

```php
class ClassName
abstract class ClassName
final class ClassName
interface InterfaceName
trait TraitName
```

This ensures non-standard files are still indexed and discoverable as dependencies.

---

## Excluded Directories

The following directories are automatically skipped during scanning:

- `vendor/` — Composer dependencies
- `node_modules/` — npm packages
- `.git/`, `.svn/` — Version control
- `.idea/` — IDE files

---

## Technical Notes

- **Binding**: Server listens on `127.0.0.1` only (not exposed to network)
- **Port**: Auto-selected free port (shown in console window)
- **Memory**: All file paths and class index stored in memory (~6K files is trivial)
- **No installation needed**: Single `.exe` file, no runtime dependencies
- **Settings persistence**: Theme and font size saved to browser `localStorage`

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Browser doesn't open | Check the console window for the URL, open it manually |
| Folder dialog doesn't appear | Make sure PowerShell is available on your system |
| Missing dependencies | Try enabling "Parse require/include" for non-autoloaded files |
| Wrong class mapping | Check Settings → Framework tab, adjust prefix mappings |
| Files not found after scan | Check that the project path is correct and .php files exist |
