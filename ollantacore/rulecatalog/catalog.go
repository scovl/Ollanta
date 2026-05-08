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
	{Key: "typescript", Name: "TypeScript", HasParser: true, HasRules: true},
	{Key: "python", Name: "Python", HasParser: true, HasRules: true},
	{Key: "rust", Name: "Rust", HasParser: true, ParserOnly: true},
}

var builtinRules = []*domain.Rule{
	// ── Go rules (14) ────────────────────────────────────────────────────────
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

	ruleDetail("go:useless-eqeq", "Useless Self-Comparison", "go", domain.TypeBug, domain.SeverityMinor,
		"Comparing a variable to itself (x == x or x != x) is always deterministic and indicates a bug or dead code.",
		"func check(a, b int) bool {\n    return a == a\n}",
		"func check(a, b int) bool {\n    return a == b\n}",
		[]string{"correctness"}),
	ruleDetail("go:useless-ifelse", "Useless If/Else", "go", domain.TypeCodeSmell, domain.SeverityMinor,
		"If statements with constant true or false conditions are dead code and should be removed or corrected.",
		"func process(x int) int {\n    if true {\n        return x * 2\n    }\n    return x\n}",
		"func process(x int) int {\n    return x * 2\n}",
		[]string{"correctness", "dead-code"}),
	ruleDetail("go:use-filepath-join", "Use filepath.Join for Path Construction", "go", domain.TypeCodeSmell, domain.SeverityMinor,
		"Building file paths by concatenating strings with + or fmt.Sprintf is fragile and platform-dependent. Use filepath.Join instead.",
		"path := dir + \"/\" + filename\nfull := fmt.Sprintf(\"%s/%s\", base, name)",
		"path := filepath.Join(dir, filename)\nfull := filepath.Join(base, name)",
		[]string{"correctness", "portability"}),
	ruleDetail("go:bad-tmp", "Insecure Temporary File", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"Creating temporary files with hardcoded /tmp/ paths is insecure. Use os.CreateTemp to ensure atomic creation and safe permissions.",
		"f, err := os.Create(\"/tmp/myapp-data.txt\")\nif err != nil {\n    log.Fatal(err)\n}",
		"f, err := os.CreateTemp(\"\", \"myapp-data-*.txt\")\nif err != nil {\n    log.Fatal(err)\n}\ndefer os.Remove(f.Name())",
		[]string{"security", "cwe-377"}),
	ruleDetail("go:math-random", "Weak Random Number Generator", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"math/rand is not cryptographically secure and must not be used for tokens, passwords, session IDs, or other security-sensitive values.",
		"import \"math/rand\"\n\ntoken := rand.Intn(1000000)\nsessionID := rand.Uint32()",
		"import \"crypto/rand\"\n\nfunc randomToken() (string, error) {\n    b := make([]byte, 16)\n    if _, err := rand.Read(b); err != nil {\n        return \"\", err\n    }\n    return hex.EncodeToString(b), nil\n}",
		[]string{"security", "cwe-338"}),
	ruleDetail("go:md5-used-as-password", "MD5 Used for Password Hashing", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"MD5 is cryptographically broken and must not be used for password hashing or any security-sensitive checksum.",
		"import \"crypto/md5\"\n\nhash := md5.Sum([]byte(password))\nstore(hex.EncodeToString(hash[:]))",
		"import \"golang.org/x/crypto/bcrypt\"\n\nhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)\nif err != nil {\n    return err\n}\nstore(string(hash))",
		[]string{"security", "cwe-916"}),

	// ── JavaScript / TypeScript rules (10) ────────────────────────────────────
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

	ruleDetail("js:useless-eqeq", "Useless Self-Comparison (JavaScript)", "javascript", domain.TypeBug, domain.SeverityMinor,
		"Comparing a variable to itself (x == x or x != x) is always deterministic and indicates a bug or dead code.",
		"if (a == a) {\n    return true;\n}",
		"if (a === b) {\n    return true;\n}",
		[]string{"correctness"}),
	ruleDetail("js:useless-assign", "Self-Assignment", "javascript", domain.TypeBug, domain.SeverityMinor,
		"Assigning a variable to itself (x = x) has no effect and indicates a bug.",
		"x = x;",
		"x = y;",
		[]string{"correctness"}),
	ruleDetail("js:leftover-debugging", "Leftover Debugging Statement", "javascript", domain.TypeCodeSmell, domain.SeverityMinor,
		"debugger statements and alert() calls should not be present in production code.",
		"function process(data) {\n    debugger;\n    alert('processing');\n    return transform(data);\n}",
		"function process(data) {\n    logger.debug('processing data', { id: data.id });\n    return transform(data);\n}",
		[]string{"convention", "debug"}),
	ruleDetail("js:assigned-undefined", "Assigned Undefined", "javascript", domain.TypeCodeSmell, domain.SeverityMinor,
		"Explicitly assigning undefined to a variable is redundant in JavaScript.",
		"let x = undefined;",
		"let x;",
		[]string{"convention", "readability"}),
	ruleDetail("js:detect-eval", "Dangerous eval() Use", "javascript", domain.TypeVulnerability, domain.SeverityCritical,
		"eval() executes arbitrary code and must not be used with untrusted input.",
		"const result = eval(userInput);",
		"const result = JSON.parse(userInput);",
		[]string{"security", "cwe-95"}),
	ruleDetail("ts:moment-deprecated", "moment.js Deprecated", "typescript", domain.TypeCodeSmell, domain.SeverityMinor,
		"moment.js is now a legacy project; consider modern alternatives.",
		"import moment from 'moment';\nconst d = moment();",
		"import { format } from 'date-fns';\nconst s = format(new Date(), 'yyyy-MM-dd');",
		[]string{"convention", "deprecated"}),

	// ── Python rules (13) ─────────────────────────────────────────────────────
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
	// Wave 1 — Python
	ruleDetail("py:useless-eqeq", "Useless Self-Comparison (Python)", "python", domain.TypeBug, domain.SeverityMinor,
		"Comparing a variable to itself (x == x or x != x) is always deterministic and indicates a bug or dead code.",
		"if a == a:\n    return True",
		"if a == b:\n    return True",
		[]string{"correctness"}),
	ruleDetail("py:useless-comparison", "Useless Comparison", "python", domain.TypeBug, domain.SeverityMinor,
		"Comparison between incompatible literal types is always deterministic and indicates a bug.",
		"if 42 == '42':\n    pass",
		"if str(42) == '42':\n    pass",
		[]string{"correctness"}),
	ruleDetail("py:dict-modify-iterating", "Dictionary Modified During Iteration", "python", domain.TypeBug, domain.SeverityMajor,
		"Deleting items from a dictionary while iterating over it raises RuntimeError.",
		"for k in d.keys():\n    if k.startswith('_'):\n        del d[k]",
		"for k in list(d.keys()):\n    if k.startswith('_'):\n        del d[k]",
		[]string{"correctness", "runtime-error"}),
	ruleDetail("py:list-modify-iterating", "List Modified During Iteration", "python", domain.TypeBug, domain.SeverityMajor,
		"Modifying a list (remove, pop, append) while iterating over it produces unexpected results.",
		"for item in items:\n    if item.is_expired():\n        items.remove(item)",
		"items = [item for item in items if not item.is_expired()]",
		[]string{"correctness", "runtime-error"}),
	ruleDetail("py:return-in-init", "Return in __init__", "python", domain.TypeBug, domain.SeverityMajor,
		"__init__ methods should not contain explicit return statements.",
		"class User:\n    def __init__(self, name):\n        if not name:\n            return\n        self.name = name",
		"class User:\n    def __init__(self, name):\n        if not name:\n            raise ValueError('name required')\n        self.name = name",
		[]string{"correctness"}),
	ruleDetail("py:pass-body", "Empty Pass Body", "python", domain.TypeCodeSmell, domain.SeverityMinor,
		"Functions or classes with only a pass statement and no docstring are incomplete.",
		"def process(data):\n    pass",
		"def process(data):\n    \"\"\"Process the given data and return normalized result.\"\"\"\n    pass",
		[]string{"convention", "incomplete"}),
	ruleDetail("py:hardcoded-tmp-path", "Hardcoded Temporary Path", "python", domain.TypeCodeSmell, domain.SeverityMinor,
		"Hardcoded /tmp/ paths are not portable and may be insecure on multi-user systems.",
		"with open('/tmp/data.txt', 'w') as f:\n    f.write(data)",
		"import tempfile\nwith tempfile.NamedTemporaryFile(mode='w', delete=False) as f:\n    f.write(data)",
		[]string{"convention", "portability"}),
	ruleDetail("py:unspecified-open-encoding", "Unspecified open() Encoding", "python", domain.TypeCodeSmell, domain.SeverityMinor,
		"open() without explicit encoding may behave differently across platforms.",
		"with open('data.txt') as f:\n    content = f.read()",
		"with open('data.txt', encoding='utf-8') as f:\n    content = f.read()",
		[]string{"correctness", "portability"}),

	// ── Wave 2 — Go (6) ───────────────────────────────────────────────────────
	ruleDetail("go:bind-all", "Bind All Interfaces", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"Binding a network listener to 0.0.0.0 or :: exposes the service on every network interface, increasing the attack surface.",
		"ln, err := net.Listen(\"tcp\", \"0.0.0.0:8080\")",
		"ln, err := net.Listen(\"tcp\", \"127.0.0.1:8080\")",
		[]string{"security", "network"}),
	ruleDetail("go:missing-ssl-minversion", "Missing TLS MinVersion", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"tls.Config without an explicit MinVersion may negotiate outdated and insecure TLS versions.",
		"cfg := &tls.Config{}",
		"cfg := &tls.Config{\n  MinVersion: tls.VersionTLS12,\n}",
		[]string{"security", "tls"}),
	ruleDetail("go:weak-crypto", "Weak Cryptographic Algorithm", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"DES, RC4 and 3DES are cryptographically broken and must not be used.",
		"import \"crypto/des\"\nblock, err := des.NewCipher(key)",
		"import \"crypto/aes\"\nblock, err := aes.NewCipher(key)",
		[]string{"security", "crypto", "cwe-327"}),
	ruleDetail("go:decompression-bomb", "Decompression Bomb", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"Creating a gzip, zlib or flate reader without size limits can lead to denial of service via decompression bombs.",
		"r, err := gzip.NewReader(src)",
		"r, err := gzip.NewReader(src)\nif err != nil {\n  return err\n}\ndefer r.Close()\nlr := io.LimitReader(r, maxDecompressSize)",
		[]string{"security", "dos"}),
	ruleDetail("go:filepath-clean-misuse", "Filepath Clean Misuse", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"Using filepath.Clean as the only path sanitization before opening a file does not prevent path traversal attacks.",
		"f, err := os.Open(filepath.Clean(userPath))",
		"clean := filepath.Clean(userPath)\nif !strings.HasPrefix(clean, safeBase) {\n  return errors.New(\"path traversal detected\")\n}\nf, err := os.Open(clean)",
		[]string{"security", "path-traversal"}),
	ruleDetail("go:loop-pointer", "Loop Variable Pointer Capture", "go", domain.TypeBug, domain.SeverityMajor,
		"Capturing a loop variable by reference inside a goroutine closure causes all goroutines to see the final value of the variable.",
		"for _, item := range items {\n  go func() {\n    process(item)\n  }()\n}",
		"for _, item := range items {\n  go func(i Item) {\n    process(i)\n  }(item)\n}",
		[]string{"correctness", "concurrency"}),

	// ── Wave 2 — Python (8) ───────────────────────────────────────────────────
	ruleDetail("py:insecure-hash", "Insecure Hash Algorithm", "python", domain.TypeVulnerability, domain.SeverityMajor,
		"hashlib.md5 and hashlib.sha1 are cryptographically broken and should not be used.",
		"import hashlib\nh = hashlib.md5(data).hexdigest()",
		"import hashlib\nh = hashlib.sha256(data).hexdigest()",
		[]string{"security", "crypto", "cwe-327"}),
	ruleDetail("py:dangerous-subprocess", "Dangerous Subprocess Call", "python", domain.TypeVulnerability, domain.SeverityMajor,
		"Using shell=True in subprocess calls allows shell injection if any argument is untrusted.",
		"subprocess.run(f'ls {user_dir}', shell=True)",
		"subprocess.run(['ls', user_dir])",
		[]string{"security", "command-injection", "cwe-78"}),
	ruleDetail("py:dangerous-os-exec", "Dangerous os.exec Call", "python", domain.TypeVulnerability, domain.SeverityMajor,
		"os.system and os.exec* functions execute arbitrary commands and are dangerous with untrusted input.",
		"os.system(f'rm {user_path}')",
		"import subprocess\nsubprocess.run(['rm', user_path])",
		[]string{"security", "command-injection", "cwe-78"}),
	ruleDetail("py:sync-sleep-in-async", "Synchronous Sleep in Async Function", "python", domain.TypeBug, domain.SeverityMajor,
		"time.sleep blocks the entire event loop in async code; use asyncio.sleep instead.",
		"async def fetch():\n    time.sleep(1)\n    return await get()",
		"async def fetch():\n    await asyncio.sleep(1)\n    return await get()",
		[]string{"correctness", "async", "performance"}),
	ruleDetail("py:missing-hash-with-eq", "Missing __hash__ with __eq__", "python", domain.TypeCodeSmell, domain.SeverityMinor,
		"A class that defines __eq__ but not __hash__ becomes unhashable in Python 3.",
		"class Point:\n    def __eq__(self, other):\n        return self.x == other.x",
		"class Point:\n    def __eq__(self, other):\n        return self.x == other.x\n    def __hash__(self):\n        return hash(self.x)",
		[]string{"correctness", "pythonic"}),
	ruleDetail("py:unchecked-returns", "Unchecked Return Value", "python", domain.TypeBug, domain.SeverityMinor,
		"Return value of os functions that can fail is discarded without error checking.",
		"os.remove(path)\nos.rename(old, new)",
		"import os\ntry:\n    os.remove(path)\nexcept OSError as e:\n    log.error(e)",
		[]string{"correctness", "error-handling"}),
	ruleDetail("py:open-never-closed", "Open Never Closed", "python", domain.TypeCodeSmell, domain.SeverityMinor,
		"open() is called but the file handle is discarded and never closed.",
		"open('data.txt')",
		"with open('data.txt') as f:\n    data = f.read()",
		[]string{"correctness", "resource-leak"}),
	ruleDetail("py:use-defused-xml", "Use defusedxml", "python", domain.TypeVulnerability, domain.SeverityMajor,
		"xml.etree.ElementTree is vulnerable to XML bombs and external entity attacks.",
		"import xml.etree.ElementTree as ET\nET.parse(user_xml)",
		"import defusedxml.ElementTree as ET\nET.parse(user_xml)",
		[]string{"security", "xml", "cwe-611"}),

	// ── Wave 2 — JavaScript / TypeScript (4) ──────────────────────────────────
	ruleDetail("js:detect-child-process", "Dangerous child_process Usage", "javascript", domain.TypeVulnerability, domain.SeverityMajor,
		"exec, execSync, spawn and spawnSync from child_process can execute arbitrary shell commands.",
		"const cp = require('child_process');\ncp.exec(userInput);",
		"const { execFile } = require('child_process');\nexecFile('ls', [userDir]);",
		[]string{"security", "command-injection", "cwe-78"}),
	ruleDetail("js:detect-insecure-websocket", "Insecure WebSocket", "javascript", domain.TypeVulnerability, domain.SeverityMajor,
		"WebSocket connections over ws:// are unencrypted and can be intercepted.",
		"const ws = new WebSocket('ws://example.com/socket');",
		"const ws = new WebSocket('wss://example.com/socket');",
		[]string{"security", "network", "cwe-319"}),
	ruleDetail("js:detect-pseudoRandomBytes", "Insecure Random Bytes", "javascript", domain.TypeVulnerability, domain.SeverityMajor,
		"crypto.pseudoRandomBytes is not cryptographically secure and must not be used for secrets.",
		"const buf = crypto.pseudoRandomBytes(16);",
		"const buf = crypto.randomBytes(16);",
		[]string{"security", "crypto", "cwe-338"}),
	ruleDetail("ts:useless-ternary", "Useless Ternary", "typescript", domain.TypeBug, domain.SeverityMinor,
		"A ternary expression that returns true/false based on the condition is redundant.",
		"const ok = result ? true : false;",
		"const ok = !!result;",
		[]string{"correctness", "readability"}),

	// ── Wave 3 — Go (5) ───────────────────────────────────────────────────────
	ruleDetail("go:cookie-missing-httponly", "Cookie Missing HttpOnly", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"http.Cookie without HttpOnly=true is vulnerable to XSS attacks via JavaScript access.",
		"&http.Cookie{Name: \"session\", Value: token}",
		"&http.Cookie{Name: \"session\", Value: token, HttpOnly: true}",
		[]string{"security", "xss", "cwe-1004"}),
	ruleDetail("go:cookie-missing-secure", "Cookie Missing Secure", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"http.Cookie without Secure=true may be transmitted over unencrypted connections.",
		"&http.Cookie{Name: \"session\", Value: token}",
		"&http.Cookie{Name: \"session\", Value: token, Secure: true}",
		[]string{"security", "cwe-614"}),
	ruleDetail("go:template-html-does-not-escape", "template.HTML Does Not Escape", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"template.HTML disables automatic HTML escaping and can lead to XSS if the input is untrusted.",
		"tmpl := template.HTML(userInput)",
		"tmpl := template.HTMLEscapeString(userInput)",
		[]string{"security", "xss", "cwe-79"}),
	ruleDetail("go:unsafe", "Unsafe Package Usage", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"The unsafe package bypasses Go's type safety and memory safety guarantees.",
		"ptr := unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + offset)",
		"Use safe Go abstractions and avoid unsafe unless absolutely necessary with thorough review.",
		[]string{"security", "memory-safety"}),
	ruleDetail("go:zip", "Zip Path Traversal", "go", domain.TypeVulnerability, domain.SeverityMajor,
		"Extracting files from zip archives without validating entry names can lead to path traversal (ZipSlip).",
		"r, _ := zip.OpenReader(\"archive.zip\")\nfor _, f := range r.File {\n  rc, _ := f.Open()\n  path := filepath.Join(\"dest\", f.Name)\n  os.Create(path)\n}",
		"r, _ := zip.OpenReader(\"archive.zip\")\nfor _, f := range r.File {\n  if strings.Contains(f.Name, \"..\") {\n    continue\n  }\n  path := filepath.Join(\"dest\", filepath.Clean(f.Name))\n}",
		[]string{"security", "path-traversal", "cwe-22"}),

	// ── Wave 3 — Python (5) ───────────────────────────────────────────────────
	ruleDetail("py:avoid-pyyaml-load", "Unsafe yaml.load", "python", domain.TypeVulnerability, domain.SeverityMajor,
		"yaml.load without Loader=SafeLoader can execute arbitrary Python code from untrusted YAML.",
		"import yaml\ndata = yaml.load(stream)",
		"import yaml\ndata = yaml.safe_load(stream)\n# or yaml.load(stream, Loader=yaml.SafeLoader)",
		[]string{"security", "deserialization", "cwe-502"}),
	ruleDetail("py:pickle", "Unsafe pickle.load", "python", domain.TypeVulnerability, domain.SeverityMajor,
		"pickle.load can execute arbitrary code during deserialization; never unpickle untrusted data.",
		"import pickle\ndata = pickle.load(file)",
		"Use JSON or msgpack for untrusted data. If pickle is required, verify the data source cryptographically.",
		[]string{"security", "deserialization", "cwe-502"}),
	ruleDetail("py:marshal", "Unsafe marshal.load", "python", domain.TypeVulnerability, domain.SeverityMajor,
		"marshal.load can execute arbitrary code; do not unmarshal untrusted data.",
		"import marshal\ndata = marshal.load(file)",
		"Use JSON or msgpack for untrusted data.",
		[]string{"security", "deserialization", "cwe-502"}),
	ruleDetail("py:unverified-ssl-context", "Unverified SSL Context", "python", domain.TypeVulnerability, domain.SeverityMajor,
		"ssl._create_unverified_context disables certificate validation and is vulnerable to MITM attacks.",
		"import ssl\nctx = ssl._create_unverified_context()",
		"import ssl\nctx = ssl.create_default_context()",
		[]string{"security", "ssl", "cwe-295"}),
	ruleDetail("py:regex-dos", "Regex Denial of Service (ReDoS)", "python", domain.TypeVulnerability, domain.SeverityMajor,
		"Regex patterns with nested quantifiers can cause catastrophic backtracking and denial of service.",
		"import re\nre.match(r'(a+)+', user_input)",
		"Use a regex library with ReDoS protection or avoid nested quantifiers.",
		[]string{"security", "dos", "cwe-1333"}),

	// ── Wave 3 — JavaScript (4) ───────────────────────────────────────────────
	ruleDetail("js:detect-redos", "Regex Denial of Service (ReDoS)", "javascript", domain.TypeVulnerability, domain.SeverityMajor,
		"Regex patterns with nested quantifiers can cause catastrophic backtracking.",
		"const re = /(a+)+/;",
		"const re = /a+/;",
		[]string{"security", "dos", "cwe-1333"}),
	ruleDetail("js:path-join-resolve-traversal", "Path Traversal via path.join", "javascript", domain.TypeVulnerability, domain.SeverityMajor,
		"path.join and path.resolve with untrusted input can lead to directory traversal.",
		"const p = path.join(baseDir, req.query.file);",
		"const safe = path.normalize(req.query.file).replace(/^(\\.\\.(\\/|\\\\|$))+/, '');\nconst p = path.join(baseDir, safe);",
		[]string{"security", "path-traversal", "cwe-22"}),
	ruleDetail("js:spawn-git-clone", "Unsafe git clone Spawn", "javascript", domain.TypeVulnerability, domain.SeverityMajor,
		"Spawning git clone with untrusted arguments can lead to command injection.",
		"spawn('git', ['clone', userUrl])",
		"Validate the URL against an allowlist before passing to spawn.",
		[]string{"security", "command-injection", "cwe-78"}),
	ruleDetail("js:incomplete-sanitization", "Incomplete Sanitization", "javascript", domain.TypeVulnerability, domain.SeverityMajor,
		"Replacing only one dangerous character (e.g., <) while leaving others (e.g., >) is insufficient and can be bypassed.",
		"const clean = input.replace(/</g, '');",
		"import DOMPurify from 'dompurify';\nconst clean = DOMPurify.sanitize(input);",
		[]string{"security", "xss", "cwe-116"}),
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
