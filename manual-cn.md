# PHP Dependency Extractor - 用户手册

## 概述

PHP Dependency Extractor（PDE）是一个独立的 Windows 工具，用于从大型项目中提取 PHP 源文件及其依赖。它会基于框架命名约定和正则解析自动发现类依赖，然后将相关文件复制到独立目录，便于分析。

**使用场景**：你有一个包含 5000+ 文件的 PHP 项目，希望只提取一小部分相关文件（例如某个 Controller 及其全部依赖），用于分享给 AI 助手或进行代码评审。

---

## 快速开始

1. 双击 `php-dep-extractor.exe`
2. 浏览器会自动打开 `http://127.0.0.1:<port>`
3. 点击 Project 旁边的 **Browse** 选择 PHP 项目目录
4. 选择 **Framework**（ZF1 / CakePHP / Laravel）
5. 点击 **Scan** 扫描并建立索引
6. 在文件树中勾选要提取的文件
7. 点击 **Analyze** 分析依赖
8. 点击 **Copy Files** 导出

---

## 界面说明

### 工具栏（Toolbar）

| 控件 | 说明 |
|------|------|
| **Project** | PHP 项目根目录路径。点击 Browse 选择。 |
| **Output** | 提取文件的导出目录。若留空，会自动生成同级目录：`{project}_output_{timestamp}`。双击该输入框可恢复自动模式。 |
| **Framework** | 选择框架用于正确解析类名：ZF1、CakePHP、Laravel。 |
| **Parse require/include** | 可选开启 `require`/`include` 语句解析，结果显示在独立区域供手动勾选。 |
| **Scan** | 遍历项目目录，构建文件树和类名索引。 |
| **Analyze** | 解析已选文件中的类引用并解析依赖。 |
| **Copy Files** | 将已选文件 + 依赖文件复制到输出目录。 |
| **Settings** | 打开设置面板（主题、字体、框架前缀映射）。 |

### 左侧面板 - 文件树（File Tree）

- 显示项目中所有 `.php` 文件
- **默认折叠**，点击文件夹逐级展开
- 点击复选框（或点击文件行）进行选择/取消
- 使用搜索框按路径过滤文件，匹配目录会自动展开

### 右侧面板 - 依赖（Dependencies）

点击 Analyze 后会出现三个区域：

| 区域 | 颜色 | 说明 |
|------|------|------|
| **Selected** | 蓝色 | 你在文件树中手动勾选的文件 |
| **Dependencies** | 橙色 | 自动发现的类依赖，会显示类名、引用类型（new、extends、static 等）以及来源文件 |
| **Include/Require** | 灰色 | 仅在启用 “Parse require/include” 时显示；每条都有复选框，需要手动勾选后才会随复制导出 |

### 状态栏（Status Bar）

- 左侧：当前操作状态
- 右侧：文件计数（Selected / Dependencies / Total）和版本号

---

## 设置（Settings）

点击 **Settings** 按钮打开设置面板，共有三个标签页：

### General

- **Theme**：可选 Dark、Light、System（跟随系统外观）
- **Font Size**：11px~18px，可使用滑块或 +/- 调整，默认 13px，设置会持久化

### Framework

可配置 **ZF1 前缀映射**：把类名前缀映射到 `application/` 下的目录。

| 前缀 | 目录 | 示例 |
|------|------|------|
| `Parent_` | `parents/` | `Parent_DbTable_Car_X` → `parents/DbTable/Car/X.php` |
| `DbTable_` | `dbs/` | `DbTable_Car_CarrierCust` → `dbs/Car/CarrierCust.php` |
| `Service_` | `services/` | `Service_Sys_Breadcrumbs` → `services/Sys/Breadcrumbs.php` |
| `Model_` | `models/` | `Model_Car_CarrierCust` → `models/Car/CarrierCust.php` |
| `Form_` | `forms/` | `Form_Login` → `forms/Login.php` |

你可以新增、删除或修改行。点击 **Save Mappings** 生效。

CakePHP 和 Laravel 标签页会展示对应检测方式说明。

### About

显示版本信息和工具说明。

---

## 框架支持

### ZF1（Zend Framework 1）

**类名解析**：以下划线分隔的类名映射为目录路径。

```
Model_Car_CarrierCust  →  application/models/Car/CarrierCust.php
DbTable_Ord_Order      →  application/dbs/Ord/Order.php
Service_Api_StarTrack  →  application/services/Api/StarTrack.php
```

**检测模式**：
- `new ClassName()`
- `extends ClassName`
- `implements InterfaceName`
- `ClassName::method()`（静态调用）
- `function foo(ClassName $bar)`（类型提示）

**会被排除**：以 `Zend_`、`ZendX_` 开头的类，以及 PHP 内置类型（`stdClass`、`Exception`、`DateTime` 等）

### CakePHP

**类名解析**：基于 CakePHP 2.x 目录约定。

| 类型 | 目录 |
|------|------|
| Model | `Model/` |
| Controller | `Controller/` |
| Component | `Controller/Component/` |
| Behavior | `Model/Behavior/` |
| Helper | `View/Helper/` |

**检测模式**：
- `App::uses('ClassName', 'Type')`
- `App::import('Type', 'ClassName')`

### Laravel

**类名解析**：PSR-4 命名空间到路径映射。

```
App\Models\User                           →  app/Models/User.php
App\Http\Controllers\OrderController     →  app/Http/Controllers/OrderController.php
App\Services\ShippingService              →  app/Services/ShippingService.php
```

**检测模式**：
- `use App\Models\User` 语句
- 带命名空间的类引用

---

## Require/Include 解析

当勾选该选项后，PDE 会扫描已选文件中的：

- `require 'path'`
- `require_once 'path'`
- `include 'path'`
- `include_once 'path'`

**支持的路径模式**：

| 模式 | 示例 |
|------|------|
| 直接字符串 | `require_once 'lib/helper.php'` |
| APPLICATION_PATH | `require_once APPLICATION_PATH . '/configs/constants.php'` |
| dirname(__FILE__) | `include dirname(__FILE__) . '/../bootstrap.php'` |
| __DIR__ | `require __DIR__ . '/functions.php'` |

结果会显示在灰色的 “Include/Require” 区域。它们**不会自动加入导出**，请在复制前手动勾选。

---

## 输出结构

复制后的文件会保留相对项目根目录的原始结构：

```
输入项目: C:\projects\myapp\
已选文件: application/controllers/V3/CustomersController.php
依赖文件: application/services/Sys/Breadcrumbs.php
         application/dbs/Car/CarrierCust.php

输出目录: C:\projects\myapp_output_20260223_153045\
        └── application/
            ├── controllers/V3/CustomersController.php
            ├── services/Sys/Breadcrumbs.php
            └── dbs/Car/CarrierCust.php
```

---

## 回退类检测（Fallback）

如果文件不符合任何框架命名约定（例如独立工具类、特殊 Controller），PDE 会读取文件前 100 行，查找：

```php
class ClassName
abstract class ClassName
final class ClassName
interface InterfaceName
trait TraitName
```

这样可确保非标准文件仍可被索引并参与依赖解析。

---

## 扫描排除目录

扫描时会自动跳过以下目录：

- `vendor/`（Composer 依赖）
- `node_modules/`（npm 包）
- `.git/`、`.svn/`（版本控制目录）
- `.idea/`（IDE 配置）

---

## 技术说明

- **绑定地址**：服务仅监听 `127.0.0.1`（不对外网开放）
- **端口**：自动选择空闲端口（显示在控制台）
- **内存**：文件路径和类索引保存在内存中（约 6000 文件压力很小）
- **免安装**：单个 `.exe`，无运行时依赖
- **设置持久化**：主题与字号保存在浏览器 `localStorage`

---

## 故障排查

| 问题 | 解决方案 |
|------|----------|
| 浏览器未自动打开 | 查看控制台中的 URL，手动打开 |
| 文件夹对话框未出现 | 确认系统可用 PowerShell |
| 依赖识别不完整 | 尝试开启 “Parse require/include” 处理未走自动加载的文件 |
| 类映射不正确 | 打开 Settings → Framework，调整 prefix mappings |
| 扫描后找不到文件 | 检查项目路径是否正确，确认存在 `.php` 文件 |
