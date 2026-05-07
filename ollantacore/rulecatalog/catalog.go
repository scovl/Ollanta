// Package rulecatalog exposes the bundled Ollanta rule catalog without loading
// analyzer implementations or tree-sitter bindings.
package rulecatalog

import (
	"sort"

	"github.com/scovl/ollanta/ollantacore/domain"
)

// Language describes a scanner language and whether Ollanta ships bundled rules for it.
type Language struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	HasParser  bool   `json:"has_parser"`
	HasRules   bool   `json:"has_rules"`
	ParserOnly bool   `json:"parser_only"`
}

var supportedLanguages = []Language{
	{Key: "go", Name: "Go", HasParser: true, HasRules: true},
	{Key: "javascript", Name: "JavaScript", HasParser: true, HasRules: true},
	{Key: "typescript", Name: "TypeScript", HasParser: true, ParserOnly: true},
	{Key: "python", Name: "Python", HasParser: true, HasRules: true},
	{Key: "rust", Name: "Rust", HasParser: true, ParserOnly: true},
}

var builtinRules = []*domain.Rule{
	// ── Go rules (8) ────────────────────────────────────────────────────────
	ruleDetail("go:cognitive-complexity", "Cognitive Complexity", "go", domain.TypeCodeSmell, domain.SeverityCritical,
		"Functions with high cognitive complexity are harder to understand, test, and maintain. Cognitive complexity measures how difficult code is to understand by considering nesting depth, structural complexity, and control flow breaks.",
		"func process(items []int) int {\n    total := 0\n    for i := 0; i < len(items); i++ {\n        if items[i] > 0 {\n            if items[i]%2 == 0 {\n                if items[i] > 10 {\n                    total += items[i] * 2\n                } else if items[i] > 5 {\n                    total += items[i]\n                } else {\n                    if items[i] > 0 {\n                        total += 1\n                    }\n                }\n            }\n        }\n    }\n    return total\n}",
		"func process(items []int) int {\n    total := 0\n    for _, item := range items {\n        total += processItem(item)\n    }\n    return total\n}\n\nfunc processItem(item int) int {\n    if item <= 0 {\n        return 0\n    }\n    return calculateValue(item)\n}\n\nfunc calculateValue(item int) int {\n    if item > 10 {\n        return item * 2\n    }\n    if item > 5 {\n        return item\n    }\n    return 1\n}",
		[]string{"complexity", "readability"}, param("max_complexity", "Maximum allowed cognitive complexity", "15", "int")),
	ruleDetail("go:function-nesting-depth", "Function Nesting Depth", "go", domain.TypeCodeSmell, domain.SeverityMajor,
		"Deeply nested code is harder to follow and increases the risk of logic errors. Each additional level of nesting adds cognitive load and makes the code more difficult to test and maintain.",
		"func process(x int) {\n    if x > 0 {\n        for i := 0; i < x; i++ {\n            if i%2 == 0 {\n                switch i {\n                case 2:\n                    if x > 5 {\n                        // logically deep nest\n                    }\n                }\n            }\n        }\n    }\n}",
		"func process(x int) {\n    if x <= 0 {\n        return\n    }\n    handleEvens(x)\n}\n\nfunc handleEvens(limit int) {\n    for i := 0; i < limit; i++ {\n        if i%2 != 0 {\n            continue\n        }\n        processEven(i, limit)\n    }\n}\n\nfunc processEven(val, limit int) {\n    if val == 2 && limit > 5 {\n        // flat logic with guard clauses\n    }\n}",
		[]string{"complexity", "readability"}, param("max_depth", "Maximum allowed nesting depth", "4", "int")),
	ruleDetail("go:magic-number", "Magic Number", "go", domain.TypeCodeSmell, domain.SeverityMinor,
		"Magic numbers are literal values used directly in code without explanation. They make code harder to understand and maintain because the meaning of the value is not clear. Extracting them into named constants improves readability.",
		"func calculatePrice(base float64) float64 {\n    return base * 1.21\n}\n\nfunc daysToSeconds(d int) int {\n    return d * 86400\n}",
		"const taxRate = 1.21\n\nfunc calculatePrice(base float64) float64 {\n    return base * taxRate\n}\n\nconst secondsPerDay = 86400\n\nfunc daysToSeconds(d int) int {\n    return d * secondsPerDay\n}",
		[]string{"readability", "convention"}, param("authorized_numbers", "Comma-separated list of allowed literal values", "0,1,2,-1", "string")),
	ruleDetail("go:naming-conventions", "Naming Conventions", "go", domain.TypeCodeSmell, domain.SeverityMinor,
		"Go has strong conventions for identifier naming. Exported identifiers should use MixedCaps, unexported identifiers should use mixedCaps. Using underscores in names goes against these conventions.",
		"type User_info struct {\n    FirstName string\n}\n\nfunc Get_User() *User_info {\n    return nil\n}",
		"type UserInfo struct {\n    FirstName string\n}\n\nfunc GetUser() *UserInfo {\n    return nil\n}",
		[]string{"convention", "readability"}),
	ruleDetail("go:no-large-functions", "No Large Functions", "go", domain.TypeCodeSmell, domain.SeverityMajor,
		"Large functions try to do too much, making them difficult to understand, test, and modify safely. Smaller functions with clear responsibilities lead to better-structured, more reusable code.",
		"func handleRequest(w http.ResponseWriter, r *http.Request) {\n    // 80+ lines of parsing, validation, business logic,\n    // database calls, error handling, and response formatting\n    // all in a single function\n}",
		"func handleRequest(w http.ResponseWriter, r *http.Request) {\n    req, err := parseRequest(r)\n    if err != nil {\n        http.Error(w, err.Error(), 400)\n        return\n    }\n    result, err := processRequest(req)\n    if err != nil {\n        http.Error(w, err.Error(), 500)\n        return\n    }\n    writeJSON(w, result)\n}",
		[]string{"size", "complexity"}, param("max_lines", "Maximum allowed lines per function", "40", "int")),
	ruleDetail("go:no-naked-returns", "No Naked Returns", "go", domain.TypeBug, domain.SeverityCritical,
		"Naked returns in Go functions with named return values reduce readability, especially in longer functions where the returned values are not immediately visible to the reader.",
		"func split(sum int) (x, y int) {\n    x = sum * 4 / 9\n    y = sum - x\n    return\n}",
		"func split(sum int) (x, y int) {\n    x = sum * 4 / 9\n    y = sum - x\n    return x, y\n}",
		[]string{"correctness", "readability"}, param("min_lines", "Minimum function length to flag naked returns", "5", "int")),
	ruleDetail("go:todo-comment", "TODO Comment", "go", domain.TypeCodeSmell, domain.SeverityInfo,
		"TODO comments indicate incomplete work and can accumulate over time, creating technical debt that is never addressed. Tracking them helps teams manage and reduce this debt.",
		"func calculate(x int) int {\n    // TODO: handle edge case\n    return x * 2\n}",
		"func calculate(x int) int {\n    // Edge case handled by validateInput called before this function.\n    return x * 2\n}",
		[]string{"convention"}),
	ruleDetail("go:too-many-parameters", "Too Many Parameters", "go", domain.TypeCodeSmell, domain.SeverityMajor,
		"Functions with too many parameters are hard to call correctly, easy to misorder, and difficult to extend. Using a parameter struct or splitting the function into smaller functions improves clarity.",
		"func CreateUser(name, email, phone, addr, city, state, zip, country string) error {\n    // ...\n}",
		"type CreateUserRequest struct {\n    Name    string\n    Email   string\n    Phone   string\n    Address string\n    City    string\n    State   string\n    Zip     string\n    Country string\n}\n\nfunc CreateUser(req CreateUserRequest) error {\n    // ...\n}",
		[]string{"design", "readability"}, param("max_params", "Maximum allowed parameter count", "5", "int")),

	// ── JavaScript rules (4) ─────────────────────────────────────────────────
	ruleDetail("js:eqeqeq", "Strict Equality (JavaScript)", "javascript", domain.TypeBug, domain.SeverityMajor,
		"Using == or != in JavaScript can produce unexpected results due to type coercion. Always use === and !== for predictable comparisons.",
		"if (value == null) {\n    return 0;\n}\n\nif (items.length != 0) {\n    process(items);\n}",
		"if (value === null) {\n    return 0;\n}\n\nif (items.length !== 0) {\n    process(items);\n}",
		[]string{"correctness", "pitfall"}),
	ruleDetail("js:no-console-log", "No console.log (JavaScript)", "javascript", domain.TypeCodeSmell, domain.SeverityMinor,
		"console.log statements left in production code clutter the console output and may leak sensitive information. Use a proper logging framework for production diagnostics.",
		"function processData(data) {\n    console.log('processing:', data);\n    return transform(data);\n}",
		"function processData(data) {\n    logger.debug('processing data', { id: data.id });\n    return transform(data);\n}",
		[]string{"convention", "debug"}),
	ruleDetail("js:no-large-functions", "No Large Functions (JS)", "javascript", domain.TypeCodeSmell, domain.SeverityMajor,
		"Large functions try to do too much, making them difficult to understand, test, and modify safely. Breaking them into smaller, focused functions improves maintainability.",
		"function handleFormSubmit(e) {\n    e.preventDefault();\n    const data = new FormData(e.target);\n    // 50+ lines of validation, API calls,\n    // state updates, and DOM manipulation\n    // all inside one function\n}",
		"function handleFormSubmit(e) {\n    e.preventDefault();\n    const data = collectFormData(e.target);\n    const errors = validateForm(data);\n    if (errors.length > 0) {\n        showErrors(errors);\n        return;\n    }\n    submitData(data);\n}",
		[]string{"size", "complexity"}, param("max_lines", "Maximum allowed lines per function", "40", "int")),
	ruleDetail("js:too-many-parameters", "Too Many Parameters (JavaScript)", "javascript", domain.TypeCodeSmell, domain.SeverityMajor,
		"Functions with too many parameters are hard to call correctly and easy to misorder. Use an options object to group related parameters together.",
		"function createUser(name, email, phone, address, city, state, zip, country) {\n    // ...\n}",
		"function createUser({ name, email, phone, address, city, state, zip, country }) {\n    // ...\n}",
		[]string{"design", "readability"}, param("max_params", "Maximum allowed parameter count", "5", "int")),

	// ── Python rules (5) ─────────────────────────────────────────────────────
	ruleDetail("py:broad-except", "Broad Exception Catch (Python)", "python", domain.TypeBug, domain.SeverityMajor,
		"Catching a broad Exception or BaseException can hide unexpected errors and make debugging difficult. Catch specific exceptions to handle only the errors you expect.",
		"try:\n    process(data)\nexcept Exception:\n    pass\n\ntry:\n    save(record)\nexcept:\n    log.error('failed')",
		"try:\n    process(data)\nexcept ValueError as e:\n    log.warning('invalid data: %s', e)\nexcept IOError as e:\n    log.error('io error: %s', e)",
		[]string{"error-handling", "correctness"}),
	ruleDetail("py:comparison-to-none", "Comparison to None (Python)", "python", domain.TypeCodeSmell, domain.SeverityMinor,
		"In Python, comparisons to None should use is or is not, not == or !=. None is a singleton and identity comparison is both more correct and more performant.",
		"if value == None:\n    return default\n\nif result != None:\n    process(result)",
		"if value is None:\n    return default\n\nif result is not None:\n    process(result)",
		[]string{"convention", "correctness"}),
	ruleDetail("py:mutable-default-argument", "Mutable Default Argument (Python)", "python", domain.TypeBug, domain.SeverityMajor,
		"Default argument values in Python are evaluated once at function definition time, not each call. Using mutable objects like lists or dicts as defaults can cause unexpected shared state across calls.",
		"def add_item(item, items=[]):\n    items.append(item)\n    return items\n\ndef configure(opts, defaults={}):\n    defaults.update(opts)\n    return defaults",
		"def add_item(item, items=None):\n    if items is None:\n        items = []\n    items.append(item)\n    return items\n\ndef configure(opts, defaults=None):\n    if defaults is None:\n        defaults = {}\n    result = {**defaults, **opts}\n    return result",
		[]string{"bug", "pitfall"}),
	ruleDetail("py:no-large-functions", "No Large Functions (Python)", "python", domain.TypeCodeSmell, domain.SeverityMajor,
		"Large functions try to do too much, making them difficult to understand, test, and modify safely. Breaking them into smaller, focused functions improves maintainability.",
		"def handle_request(request):\n    # 60+ lines of parsing, validation, business\n    # logic, database queries, and response formatting\n    # all inside one function\n    return response",
		"def handle_request(request):\n    data = parse_request(request)\n    errors = validate(data)\n    if errors:\n        return error_response(errors)\n    result = process(data)\n    return success_response(result)",
		[]string{"size", "complexity"}, param("max_lines", "Maximum allowed lines per function", "40", "int")),
	ruleDetail("py:too-many-parameters", "Too Many Parameters (Python)", "python", domain.TypeCodeSmell, domain.SeverityMajor,
		"Functions with too many parameters are hard to call correctly and easy to misorder. Use a dataclass, TypedDict, or keyword-only arguments to improve clarity.",
		"def create_user(name, email, phone, address, city, state, zip_code, country):\n    pass",
		"from dataclasses import dataclass\n\n@dataclass\nclass CreateUserRequest:\n    name: str\n    email: str\n    phone: str = ''\n    address: str = ''\n    city: str = ''\n    state: str = ''\n    zip_code: str = ''\n    country: str = ''\n\ndef create_user(req: CreateUserRequest):\n    pass",
		[]string{"design", "readability"}, param("max_params", "Maximum allowed parameter count (excluding self/cls)", "5", "int")),
}

// SupportedLanguages returns all languages known to the scanner.
func SupportedLanguages() []Language {
	out := make([]Language, len(supportedLanguages))
	copy(out, supportedLanguages)
	return out
}

// Rules returns all bundled rule metadata as defensive copies.
func Rules() []*domain.Rule {
	out := make([]*domain.Rule, 0, len(builtinRules))
	for _, rule := range builtinRules {
		out = append(out, cloneRule(rule))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// ByKey returns a bundled rule by key.
func ByKey(key string) (*domain.Rule, bool) {
	for _, rule := range builtinRules {
		if rule.Key == key {
			return cloneRule(rule), true
		}
	}
	return nil, false
}

// ByLanguage returns bundled rules for the given language.
func ByLanguage(language string) []*domain.Rule {
	out := []*domain.Rule{}
	for _, rule := range builtinRules {
		if rule.Language == language || rule.Language == "*" {
			out = append(out, cloneRule(rule))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// LanguageByKey returns metadata for a supported language.
func LanguageByKey(key string) (Language, bool) {
	for _, language := range supportedLanguages {
		if language.Key == key {
			return language, true
		}
	}
	return Language{}, false
}

// LanguageHasRules reports whether the language has at least one bundled rule.
func LanguageHasRules(key string) bool {
	language, ok := LanguageByKey(key)
	return ok && language.HasRules
}

// LanguageIsParserOnly reports whether Ollanta can parse the language but ships no rules yet.
func LanguageIsParserOnly(key string) bool {
	language, ok := LanguageByKey(key)
	return ok && language.ParserOnly
}

// DefaultParams returns rule parameter defaults keyed by parameter name.
func DefaultParams(rule *domain.Rule) map[string]string {
	out := make(map[string]string, len(rule.ParamsSchema))
	for key, param := range rule.ParamsSchema {
		out[key] = param.DefaultValue
	}
	return out
}

func rule(key, name, language string, issueType domain.IssueType, severity domain.Severity, tags []string, params ...domain.ParamDef) *domain.Rule {
	schema := make(map[string]domain.ParamDef, len(params))
	for _, p := range params {
		schema[p.Key] = p
	}
	return &domain.Rule{
		Key:             key,
		Name:            name,
		Language:        language,
		Type:            issueType,
		DefaultSeverity: severity,
		Tags:            append([]string(nil), tags...),
		ParamsSchema:    schema,
	}
}

func ruleDetail(key, name, language string, issueType domain.IssueType, severity domain.Severity, rationale, noncompliantCode, compliantCode string, tags []string, params ...domain.ParamDef) *domain.Rule {
	r := rule(key, name, language, issueType, severity, tags, params...)
	r.Rationale = rationale
	r.NoncompliantCode = noncompliantCode
	r.CompliantCode = compliantCode
	return r
}

func param(key, description, defaultValue, paramType string) domain.ParamDef {
	return domain.ParamDef{Key: key, Description: description, DefaultValue: defaultValue, Type: paramType}
}

func cloneRule(rule *domain.Rule) *domain.Rule {
	clone := *rule
	clone.Tags = append([]string(nil), rule.Tags...)
	clone.ParamsSchema = make(map[string]domain.ParamDef, len(rule.ParamsSchema))
	for key, param := range rule.ParamsSchema {
		clone.ParamsSchema[key] = param
	}
	return &clone
}
