package parser

import (
	"os"
	"regexp"
	"strings"
)

// ClassReference represents a found class reference in a PHP file.
type ClassReference struct {
	ClassName string `json:"className"`
	RefType   string `json:"refType"` // "new", "extends", "implements", "static", "typehint", "use"
	Line      int    `json:"line"`
}

// PHP built-in types/classes to exclude.
var builtinClasses = map[string]bool{
	"self": true, "static": true, "parent": true,
	"stdClass": true, "Exception": true, "RuntimeException": true,
	"InvalidArgumentException": true, "LogicException": true,
	"DateTime": true, "DateTimeImmutable": true, "DateInterval": true,
	"ArrayObject": true, "ArrayIterator": true, "Iterator": true,
	"Countable": true, "Serializable": true, "JsonSerializable": true,
	"Closure": true, "Generator": true, "Throwable": true, "Error": true,
	"TypeError": true, "ValueError": true,
	"PDO": true, "PDOStatement": true, "PDOException": true,
	"SplFileInfo": true, "SplFileObject": true, "SplHeap": true,
	"SplStack": true, "SplQueue": true, "SplPriorityQueue": true,
	"DOMDocument": true, "DOMElement": true, "DOMNode": true,
	"SimpleXMLElement": true, "XMLReader": true, "XMLWriter": true,
	"ReflectionClass": true, "ReflectionMethod": true,
	"SoapClient": true, "SoapServer": true,
	"mysqli": true, "mysqli_result": true,
	"null": true, "true": true, "false": true,
	"int": true, "float": true, "string": true, "bool": true,
	"array": true, "object": true, "void": true, "mixed": true,
	"callable": true, "iterable": true, "never": true,
}

// Framework prefixes to exclude.
var frameworkPrefixes = []string{
	"Zend_", "ZendX_",
	"Cake", // CakePHP core
	"Illuminate\\",
	"Symfony\\",
	"PHPUnit",
}

// Regex patterns for class references.
var (
	reNew        = regexp.MustCompile(`new\s+([A-Z]\w+)`)
	reExtends    = regexp.MustCompile(`extends\s+([A-Z]\w+)`)
	reImplements = regexp.MustCompile(`implements\s+(.+?)[\s{]`)
	reStatic     = regexp.MustCompile(`([A-Z]\w+)::`)
	reTypeHint   = regexp.MustCompile(`(?:function\s+\w+\s*\([^)]*?)([A-Z]\w+)\s+\$`)
	reUseStmt    = regexp.MustCompile(`^use\s+(App\\[\w\\]+)`)
	// ZF1 style: class names with underscores like Model_Car_CarrierCust
	reZF1Class = regexp.MustCompile(`new\s+([A-Z]\w*(?:_\w+)+)`)
	// CakePHP App::uses
	reCakeUses  = regexp.MustCompile(`App::uses\s*\(\s*'(\w+)'`)
	reCakeImport = regexp.MustCompile(`App::import\s*\(\s*'(\w+)'\s*,\s*'(\w+)'`)
)

// ExtractClassRefs extracts all class references from a PHP file.
func ExtractClassRefs(filePath string) ([]ClassReference, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	var refs []ClassReference
	seen := make(map[string]bool)

	addRef := func(className, refType string, line int) {
		className = strings.TrimSpace(className)
		if className == "" || builtinClasses[className] {
			return
		}
		if isFrameworkClass(className) {
			return
		}
		key := className + "|" + refType
		if !seen[key] {
			seen[key] = true
			refs = append(refs, ClassReference{
				ClassName: className,
				RefType:   refType,
				Line:      line,
			})
		}
	}

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Skip comments
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "/*") {
			continue
		}

		// new ClassName
		for _, m := range reNew.FindAllStringSubmatch(line, -1) {
			addRef(m[1], "new", lineNum)
		}
		// ZF1 style new with underscores
		for _, m := range reZF1Class.FindAllStringSubmatch(line, -1) {
			addRef(m[1], "new", lineNum)
		}

		// extends ClassName
		if m := reExtends.FindStringSubmatch(line); len(m) > 1 {
			addRef(m[1], "extends", lineNum)
		}

		// implements Interface1, Interface2
		if m := reImplements.FindStringSubmatch(line); len(m) > 1 {
			for _, iface := range strings.Split(m[1], ",") {
				addRef(strings.TrimSpace(iface), "implements", lineNum)
			}
		}

		// ClassName::method()
		for _, m := range reStatic.FindAllStringSubmatch(line, -1) {
			name := m[1]
			if name != "self" && name != "static" && name != "parent" {
				addRef(name, "static", lineNum)
			}
		}

		// Type hints in function params
		for _, m := range reTypeHint.FindAllStringSubmatch(line, -1) {
			addRef(m[1], "typehint", lineNum)
		}

		// Laravel use statements
		if m := reUseStmt.FindStringSubmatch(trimmed); len(m) > 1 {
			addRef(m[1], "use", lineNum)
		}

		// CakePHP App::uses
		for _, m := range reCakeUses.FindAllStringSubmatch(line, -1) {
			addRef(m[1], "uses", lineNum)
		}
		for _, m := range reCakeImport.FindAllStringSubmatch(line, -1) {
			addRef(m[2], "import", lineNum)
		}
	}

	return refs, nil
}

func isFrameworkClass(name string) bool {
	for _, prefix := range frameworkPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
