package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var supportedLoggers = map[string]bool{
	"log/slog": true,
	"zap":      true,
	"slog":     true,
}

var loggerMethods = map[string]bool{
	"Info":    true,
	"Error":   true,
	"Warn":    true,
	"Debug":   true,
	"Print":   true,
	"Printf":  true,
	"Println": true,
}

var sensitiveKeywords = []string{
	"password",
	"passwd",
	"pwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"key",
	"auth",
	"credential",
	"private",
}

var Analyzer = &analysis.Analyzer{
	Name: "logrules",
	Doc:  "проверяет что лог-сообщения соответствуют правилам: строчная буква, только английский, без спецсимволов, без чувствительных данных",
	Run:  run,
}

type loggerInfo struct {
	pkgPath string
	method  string
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		imports := getImports(file)

		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			loggerInfo := isLoggerCall(callExpr, pass, imports)
			if loggerInfo == nil {
				return true
			}

			message := extractMessage(callExpr)
			if message == "" {
				return true
			}

			checkRules(message, callExpr, pass)

			return true
		})
	}

	return nil, nil
}

func isLoggerCall(call *ast.CallExpr, pass *analysis.Pass, imports map[string]string) *loggerInfo {
	selExpr, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	methodName := selExpr.Sel.Name
	if !loggerMethods[methodName] {
		return nil
	}

	switch x := selExpr.X.(type) {
	case *ast.Ident:

		pkgName := x.Name

		if pkgPath, ok := imports[pkgName]; ok {
			for supported := range supportedLoggers {
				if strings.Contains(pkgPath, supported) {
					return &loggerInfo{
						pkgPath: pkgPath,
						method:  methodName,
					}
				}
			}
		}

		obj := pass.TypesInfo.ObjectOf(x)
		if obj != nil {
			typeStr := obj.Type().String()
			for supported := range supportedLoggers {
				if strings.Contains(typeStr, supported) {
					return &loggerInfo{
						pkgPath: typeStr,
						method:  methodName,
					}
				}
			}
		}

	case *ast.SelectorExpr:

		if pkgIdent, ok := x.X.(*ast.Ident); ok {
			if pkgPath, ok := imports[pkgIdent.Name]; ok {
				for supported := range supportedLoggers {
					if strings.Contains(pkgPath, supported) {
						return &loggerInfo{
							pkgPath: pkgPath,
							method:  methodName,
						}
					}
				}
			}
		}
	}

	return nil
}

func getImports(file *ast.File) map[string]string {
	imports := make(map[string]string)

	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, "\"")

		var name string
		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			parts := strings.Split(path, "/")
			name = parts[len(parts)-1]
		}

		imports[name] = path
	}

	return imports
}

func extractMessage(call *ast.CallExpr) string {
	if len(call.Args) == 0 {
		return ""
	}

	firstArg := call.Args[0]

	switch arg := firstArg.(type) {
	case *ast.BasicLit:
		if arg.Kind == token.STRING {
			str := strings.Trim(arg.Value, "\"")
			str = strings.Trim(str, "`")
			return str
		}
	case *ast.Ident:
		return arg.Name + "_variable_reference"
	case *ast.BinaryExpr:
		if arg.Op == token.ADD {
			left := ""
			if leftLit, ok := arg.X.(*ast.BasicLit); ok && leftLit.Kind == token.STRING {
				left = strings.Trim(leftLit.Value, "\"")
			}
			right := ""
			if rightIdent, ok := arg.Y.(*ast.Ident); ok {
				right = rightIdent.Name
			}

			if left != "" && right != "" {
				return left + "_concat_" + right
			}
		}
		return ""
	}
	return ""
}

func checkRules(message string, call *ast.CallExpr, pass *analysis.Pass) {
	checkLowercase(message, call, pass)
	checkEnglish(message, call, pass)
	checkSpecialChars(message, call, pass)
	checkSensitive(message, call, pass)
}

func checkLowercase(message string, call *ast.CallExpr, pass *analysis.Pass) {
	if message == "" {
		return
	}

	runes := []rune(message)
	if len(runes) == 0 {
		return
	}

	firstChar := runes[0]

	if unicode.IsLetter(firstChar) && !unicode.IsLower(firstChar) {
		pass.Reportf(call.Pos(), "лог-сообщение должно начинаться со строчной буквы (начинается с %q)", string(firstChar))
	}
}

func checkEnglish(message string, call *ast.CallExpr, pass *analysis.Pass) {
	if message == "" {
		return
	}

	if strings.HasSuffix(message, "_variable_reference") {
		return
	}

	for _, r := range message {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
			pass.Reportf(call.Pos(), "лог-сообщение должно содержать только английские буквы (найден символ %q)", string(r))
			return
		}
	}
}

func checkSpecialChars(message string, call *ast.CallExpr, pass *analysis.Pass) {
	if message == "" {
		return
	}

	if strings.HasSuffix(message, "_variable_reference") {
		return
	}

	for _, r := range message {
		if isEmoji(r) {
			pass.Reportf(call.Pos(), "лог-сообщение не должно содержать эмодзи (найден символ %q)", string(r))
			return
		}
	}

	trimmed := strings.TrimRight(message, " ")
	if strings.HasSuffix(trimmed, "!") ||
		strings.HasSuffix(trimmed, "?") ||
		strings.HasSuffix(trimmed, "...") {
		pass.Reportf(call.Pos(), "лог-сообщение не должно заканчиваться на ! ? или ...")
		return
	}

	if strings.Contains(message, "!!") ||
		strings.Contains(message, "??") ||
		strings.Contains(message, "...") {
		pass.Reportf(call.Pos(), "лог-сообщение не должно содержать повторяющиеся знаки препинания")
		return
	}
}

func checkSensitive(message string, call *ast.CallExpr, pass *analysis.Pass) {
	if message == "" {
		return
	}

	fmt.Fprintf(os.Stderr, "    checkSensitive: проверяем сообщение %q\n", message)

	if strings.Contains(message, "_concat_") {
		fmt.Fprintf(os.Stderr, "    checkSensitive: обнаружена конкатенация\n")
		parts := strings.SplitN(message, "_concat_", 2)
		if len(parts) == 2 {
			prefix := parts[0]
			varName := parts[1]
			fmt.Fprintf(os.Stderr, "      префикс: %q, переменная: %q\n", prefix, varName)

			lowerPrefix := strings.ToLower(prefix)
			for _, keyword := range sensitiveKeywords {
				if strings.Contains(lowerPrefix, strings.ToLower(keyword)) {
					fmt.Fprintf(os.Stderr, "      найдено ключевое слово %q в префиксе\n", keyword)
					pass.Reportf(call.Pos(), "лог-сообщение может содержать чувствительные данные: найдено ключевое слово %q в префиксе", keyword)
					return
				}
			}

			lowerVar := strings.ToLower(varName)
			for _, keyword := range sensitiveKeywords {
				if strings.Contains(lowerVar, strings.ToLower(keyword)) {
					fmt.Fprintf(os.Stderr, "      найдено ключевое слово %q в имени переменной\n", keyword)
					pass.Reportf(call.Pos(), "потенциально чувствительные данные: логирование переменной %q может содержать %s", varName, keyword)
					return
				}
			}
		}
		return
	}

	if strings.HasSuffix(message, "_variable_reference") {
		varName := strings.TrimSuffix(message, "_variable_reference")
		fmt.Fprintf(os.Stderr, "    checkSensitive: переменная %q\n", varName)

		for _, keyword := range sensitiveKeywords {
			if strings.Contains(strings.ToLower(varName), strings.ToLower(keyword)) {
				fmt.Fprintf(os.Stderr, "      найдено ключевое слово %q в имени переменной\n", keyword)
				pass.Reportf(call.Pos(), "потенциально чувствительные данные: логирование переменной %q может содержать %s", varName, keyword)
				return
			}
		}
		return
	}

	lowerMessage := strings.ToLower(message)
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lowerMessage, strings.ToLower(keyword)) {
			fmt.Fprintf(os.Stderr, "    checkSensitive: найдено ключевое слово %q в тексте\n", keyword)
			pass.Reportf(call.Pos(), "лог-сообщение может содержать чувствительные данные: найдено ключевое слово %q", keyword)
			return
		}
	}
}

func isEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) ||
		(r >= 0x1F300 && r <= 0x1F5FF) ||
		(r >= 0x1F680 && r <= 0x1F6FF) ||
		(r >= 0x2600 && r <= 0x26FF) ||
		(r >= 0x2700 && r <= 0x27BF) ||
		(r >= 0xFE00 && r <= 0xFE0F) ||
		(r >= 0x1F900 && r <= 0x1F9FF)
}
