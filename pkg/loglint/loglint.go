package loglint

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func init() {
	// Регистрируем плагин для golangci-lint.
	register.Plugin("loglint", New)
}

var slogMethods = map[string]struct{}{
	"Debug": {},
	"Info":  {},
	"Warn":  {},
	"Error": {},
}

var zapLoggerMethods = map[string]struct{}{
	"Debug":  {},
	"Info":   {},
	"Warn":   {},
	"Error":  {},
	"DPanic": {},
	"Panic":  {},
	"Fatal":  {},
}

var zapSugaredMethods = map[string]struct{}{
	"Debugf":  {},
	"Infof":   {},
	"Warnf":   {},
	"Errorf":  {},
	"DPanicf": {},
	"Panicf":  {},
	"Fatalf":  {},
	"Debugw":  {},
	"Infow":   {},
	"Warnw":   {},
	"Errorw":  {},
	"DPanicw": {},
	"Panicw":  {},
	"Fatalw":  {},
}

type MySettings struct {
	RequireLowercaseStart *bool    `json:"require-lowercase-start"`
	RequireEnglish        *bool    `json:"require-english"`
	ForbidSpecialChars    *bool    `json:"forbid-special-chars"`
	ForbidSensitiveData   *bool    `json:"forbid-sensitive-data"`
	SensitiveKeywords     []string `json:"sensitive-keywords"`
	SensitivePatterns     []string `json:"sensitive-patterns"`
}

type Config struct {
	RequireLowercaseStart bool
	RequireEnglish        bool
	ForbidSpecialChars    bool
	ForbidSensitiveData   bool
	SensitiveKeywords     []string
	SensitivePatterns     []*regexp.Regexp
}

type PluginLogLint struct {
	settings Config
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[MySettings](settings)
	if err != nil {
		return nil, err
	}

	cfg, err := applyDefaults(s)
	if err != nil {
		return nil, err
	}

	return &PluginLogLint{settings: cfg}, nil
}

func (f *PluginLogLint) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name:     "loglint",
			Doc:      "Checks log message rules (lowercase, English, special chars, sensitive data)",
			Requires: []*analysis.Analyzer{inspect.Analyzer},
			Run:      f.run,
		},
	}, nil
}

func (f *PluginLogLint) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func (f *PluginLogLint) run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		selector := getSelector(call.Fun)
		if selector == nil {
			return
		}

		pkgPath, typeName, isPkgSelector := selectorTarget(pass.TypesInfo, selector)
		method := selector.Sel.Name

		if isPkgSelector && isSlogMethod(method) && pkgPath == "log/slog" {
			checkFirstArg(pass, call, f.settings)
			return
		}

		if pkgPath == "log/slog" && typeName == "Logger" && isSlogMethod(method) {
			checkFirstArg(pass, call, f.settings)
			return
		}

		if pkgPath == "go.uber.org/zap" && typeName == "Logger" && isZapMethod(method) {
			// Проверяем только первый аргумент вызова.
			checkFirstArg(pass, call, f.settings)
			return
		}

		if pkgPath == "go.uber.org/zap" && typeName == "SugaredLogger" && isZapSugaredMethod(method) {
			// Для SugaredLogger проверяем только Infof/Infow и аналоги.
			checkFirstArg(pass, call, f.settings)
		}
	})

	return nil, nil
}

func getSelector(fun ast.Expr) *ast.SelectorExpr {
	switch expr := fun.(type) {
	case *ast.SelectorExpr:
		return expr
	case *ast.IndexExpr:
		// Нужен для generic-вызовов вида logger.Info[T](...).
		return getSelector(expr.X)
	case *ast.IndexListExpr:
		// Поддержка multi-parameter generic-вызовов.
		return getSelector(expr.X)
	case *ast.ParenExpr:
		// Снимаем лишние скобки вокруг селектора.
		return getSelector(expr.X)
	default:
		return nil
	}
}

func selectorTarget(info *types.Info, sel *ast.SelectorExpr) (pkgPath, typeName string, isPkgSelector bool) {
	if ident, ok := sel.X.(*ast.Ident); ok {
		if obj, ok := info.Uses[ident].(*types.PkgName); ok {
			// Селектор пакета: slog.Info(...)
			return obj.Imported().Path(), "", true
		}
	}

	recvType := info.TypeOf(sel.X)
	if recvType == nil {
		return "", "", false
	}

	recvType = deref(recvType)
	if named, ok := recvType.(*types.Named); ok {
		obj := named.Obj()
		if obj != nil && obj.Pkg() != nil {
			// Вызов метода у типа: logger.Info(...)
			return obj.Pkg().Path(), obj.Name(), false
		}
	}

	return "", "", false
}

func deref(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		// Работаем с типом значения, а не указателя.
		return ptr.Elem()
	}
	return t
}

func isSlogMethod(name string) bool {
	_, ok := slogMethods[name]
	return ok
}

func isZapMethod(name string) bool {
	_, ok := zapLoggerMethods[name]
	return ok
}

func isZapSugaredMethod(name string) bool {
	_, ok := zapSugaredMethods[name]
	return ok
}
func applyDefaults(settings MySettings) (Config, error) {
	defaultTrue := func(v *bool) bool {
		if v == nil {
			return true
		}
		return *v
	}

	cfg := Config{
		RequireLowercaseStart: defaultTrue(settings.RequireLowercaseStart),
		RequireEnglish:        defaultTrue(settings.RequireEnglish),
		ForbidSpecialChars:    defaultTrue(settings.ForbidSpecialChars),
		ForbidSensitiveData:   defaultTrue(settings.ForbidSensitiveData),
		SensitiveKeywords:     settings.SensitiveKeywords,
	}

	if cfg.SensitiveKeywords == nil {
		cfg.SensitiveKeywords = []string{
			"password",
			"passwd",
			"pwd",
			"secret",
			"api_key",
			"api key",
			"apikey",
			"access_key",
			"private_key",
			"token",
			"bearer",
		}
	}

	for _, pattern := range settings.SensitivePatterns {
		if strings.TrimSpace(pattern) == "" {
			continue
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return Config{}, err
		}
		cfg.SensitivePatterns = append(cfg.SensitivePatterns, re)
	}

	return cfg, nil
}

func checkFirstArg(pass *analysis.Pass, call *ast.CallExpr, cfg Config) {
	if len(call.Args) == 0 {
		return
	}

	message, ok := stringConstValue(pass, call.Args[0])
	if !ok {
		return
	}
	if cfg.RequireEnglish && containsNonEnglishLetters(message) {
		pass.Reportf(call.Args[0].Pos(), "log message should contain only English letters")
	}

	if cfg.RequireLowercaseStart && !startsWithLowercase(message) {
		if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if fixed, ok := buildLowercaseFix(lit.Value, message); ok {
				pass.Report(analysis.Diagnostic{
					Pos:     call.Args[0].Pos(),
					End:     call.Args[0].End(),
					Message: "log message should start with a lowercase letter",
					SuggestedFixes: []analysis.SuggestedFix{
						{
							TextEdits: []analysis.TextEdit{
								{
									Pos:     call.Args[0].Pos(),
									End:     call.Args[0].End(),
									NewText: []byte(fixed),
								},
							},
						},
					},
				})
			} else {
				pass.Reportf(call.Args[0].Pos(), "log message should start with a lowercase letter")
			}
		} else {
			pass.Reportf(call.Args[0].Pos(), "log message should start with a lowercase letter")
		}
	}

	if cfg.ForbidSpecialChars && containsSpecialChars(message) {
		pass.Reportf(call.Args[0].Pos(), "log message should not contain special characters or emoji")
	}

	if cfg.ForbidSensitiveData && containsSensitiveData(message, cfg.SensitiveKeywords, cfg.SensitivePatterns) {
		pass.Reportf(call.Args[0].Pos(), "log message should not contain sensitive data")
	}
}

func stringConstValue(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	if tv, ok := pass.TypesInfo.Types[expr]; ok && tv.Value != nil && tv.Value.Kind() == constant.String {
		// Константные выражения берём через types.Info.
		return constant.StringVal(tv.Value), true
	}

	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		// Запасной путь: прямой строковый литерал.
		value, err := strconv.Unquote(lit.Value)
		if err == nil {
			return value, true
		}
	}

	return "", false
}

func startsWithLowercase(message string) bool {
	for _, r := range message {
		if unicode.IsSpace(r) {
			continue
		}
		return unicode.IsLower(r)
	}
	return true
}

func buildLowercaseFix(litValue, message string) (string, bool) {
	runes := []rune(message)
	for i, r := range runes {
		if unicode.IsSpace(r) {
			continue
		}
		runes[i] = unicode.ToLower(r)
		break
	}
	correctMessage := string(runes)

	if strings.HasPrefix(litValue, "`") {
		// Сохраняем raw-строку при автоправке.
		return "`" + correctMessage + "`", true
	}

	quoted := strconv.Quote(correctMessage)
	return quoted, true
}

func containsNonEnglishLetters(message string) bool {
	for _, r := range message {
		if unicode.IsLetter(r) {
			if !(r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z') {
				return true
			}
		}
	}
	return false
}

func containsSpecialChars(message string) bool {
	for _, r := range message {
		if r == ' ' {
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}
		// Намеренно строгая проверка: пунктуация и эмодзи запрещены.
		return true
	}
	return false
}

func containsSensitiveData(message string, keywords []string, patterns []*regexp.Regexp) bool {
	lower := strings.ToLower(message)
	for _, kw := range keywords {
		keyword := strings.ToLower(strings.TrimSpace(kw))
		if keyword == "" {
			continue
		}

		for idx := strings.Index(lower, keyword); idx != -1; {
			next := idx + len(keyword)
			for next < len(lower) && lower[next] == ' ' {
				next++
			}
			// Срабатываем, если после ключа есть ':' или '='.
			if next < len(lower) && (lower[next] == ':' || lower[next] == '=') {
				return true
			}
			searchFrom := idx + len(keyword)
			if searchFrom >= len(lower) {
				break
			}
			idx = strings.Index(lower[searchFrom:], keyword)
			if idx != -1 {
				idx += searchFrom
			}
		}
	}
	for _, re := range patterns {
		if re.MatchString(message) {
			return true
		}
	}
	return false
}
