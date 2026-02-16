```markdown
# Линтер для проверки лог-записей

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

Статический анализатор для Go-кода, который проверяет лог-сообщения на соответствие правилам форматирования и безопасности. Интегрируется с `go vet` и `golangci-lint`.

## 📦 Репозиторий

Проект доступен на GitHub: [https://github.com/yourusername/logger-linter](https://github.com/yourusername/logger-linter)

```bash
git clone https://github.com/yourusername/logger-linter.git
cd logger-linter
```

## 📋 Правила проверки

### 1. Строчная буква в начале
Лог-сообщения должны начинаться со строчной буквы.

```go
// ❌ Неправильно
log.Info("Server started")
slog.Error("Failed to connect")

// ✅ Правильно
log.Info("starting server")
slog.Error("failed to connect")
```

### 2. Только английский язык
Лог-сообщения должны содержать только латинские символы.

```go
// ❌ Неправильно
log.Info("запуск сервера")
slog.Error("ошибка подключения")

// ✅ Правильно
log.Info("starting server")
slog.Error("connection failed")
```

### 3. Без спецсимволов и эмодзи
Запрещены эмодзи, повторяющиеся знаки препинания, завершение на ! ? или ...

```go
// ❌ Неправильно
log.Info("server started!")
slog.Error("failed!!")
log.Warn("loading...")
slog.Info("success 😊")

// ✅ Правильно
log.Info("server started")
slog.Error("failed")
log.Warn("loading")
slog.Info("success")
```

### 4. Без чувствительных данных
Запрещено логирование паролей, токенов, ключей API.

```go
// ❌ Неправильно
log.Info("user password: " + password)
slog.Debug("api_key=" + apiKey)
log.Info("token: " + token)

// ✅ Правильно
log.Info("user authenticated")
slog.Debug("api request completed")
log.Info("token validated")
```

## 🚀 Поддерживаемые логгеры

- `log/slog` (стандартный пакет Go)
- `go.uber.org/zap`

## 🔧 Установка и сборка

### Требования
- Go 1.22 или выше
- Make (опционально)

### Сборка

```bash
# Клонирование репозитория
git clone https://github.com/yourusername/logger-linter.git
cd logger-linter

# Сборка линтера
make build
# или
go build -o mylinter ./cmd/mylinter

# Проверка сборки
./mylinter -help
```

## 📖 Использование

### Запуск через go vet

```bash
# Проверка всех пакетов проекта
go vet -vettool=./mylinter ./...

# Проверка конкретного файла
go vet -vettool=./mylinter main.go

# Проверка с подробным выводом
go vet -v -vettool=./mylinter ./...
```

### Использование Makefile

```bash
make build    # сборка линтера
make test     # запуск тестов
make vet      # проверка примеров
make clean    # очистка артефактов
make all      # полная сборка и проверка
```

## 🔌 Интеграция с golangci-lint

### Вариант 1: как плагин

1. Соберите плагин:
```bash
go build -buildmode=plugin -o mylinter.so ./cmd/mylinter
```

2. Создайте файл `.golangci.yml` в корне вашего проекта:
```yaml
version: "2"
linters:
  enable:
    - mylinter
  settings:
    custom:
      mylinter:
        type: plugin
        path: /абсолютный/путь/к/mylinter.so   # или относительный ./mylinter.so
        description: "Линтер для проверки логов"
```

3. Запустите:
```bash
golangci-lint run
```

### Вариант 2: как модуль (рекомендуется)

1. Убедитесь, что в `cmd/mylinter` есть функция `AnalyzerPlugin`, экспортирующая анализатор:
```go
package main

import (
    "github.com/yourusername/logger-linter/pkg/analyzer"
    "golang.org/x/tools/go/analysis"
)

func AnalyzerPlugin() []*analysis.Analyzer {
    return []*analysis.Analyzer{analyzer.Analyzer}
}
```

2. Добавьте replace в `go.mod` вашего проекта:
```
replace github.com/yourusername/logger-linter/cmd/mylinter => /путь/к/logger-linter/cmd/mylinter
```

3. Создайте `.golangci.yml`:
```yaml
version: "2"
linters:
  enable:
    - mylinter
  settings:
    custom:
      mylinter:
        type: module
        path: github.com/yourusername/logger-linter/cmd/mylinter
        description: "Линтер для проверки логов"
```

4. Запустите:
```bash
golangci-lint run
```

## 📝 Примеры использования

### Плохой код (с нарушениями)

```go
package main

import (
    "log/slog"
    "os"
    "go.uber.org/zap"
)

func main() {
    slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    zapLogger := zap.NewExample()

    password := "supersecret"
    apiKey := "12345"

    slogLogger.Info("Server started on port 8080")      // заглавная S
    slogLogger.Error("Failed to connect to Database")   // заглавная F
    slogLogger.Info("привет мир")                       // русские буквы
    slogLogger.Warn("warning!!!")                       // спецсимволы
    slogLogger.Debug("user password: " + password)      // чувствительные данные
    slogLogger.Info("api_key: " + apiKey)               // чувствительные данные

    zapLogger.Info("Server started")                     // заглавная S
    zapLogger.Error("ошибка подключения")                 // русские буквы
    zapLogger.Warn("loading...")                         // спецсимволы
    zapLogger.Debug("token: abc123")                     // чувствительные данные
}
```

Вывод линтера:
```
./examples/bad/main.go:17:2: лог-сообщение должно начинаться со строчной буквы (начинается с "S")
./examples/bad/main.go:18:2: лог-сообщение должно начинаться со строчной буквы (начинается с "F")
./examples/bad/main.go:19:2: лог-сообщение должно содержать только английские буквы (найден символ "п")
./examples/bad/main.go:20:2: лог-сообщение не должно заканчиваться на ! ? или ...
./examples/bad/main.go:21:2: лог-сообщение может содержать чувствительные данные: найдено ключевое слово "password" в префиксе
./examples/bad/main.go:22:2: лог-сообщение может содержать чувствительные данные: найдено ключевое слово "api_key" в префиксе
./examples/bad/main.go:24:2: лог-сообщение должно начинаться со строчной буквы (начинается с "S")
./examples/bad/main.go:25:2: лог-сообщение должно содержать только английские буквы (найден символ "о")
./examples/bad/main.go:26:2: лог-сообщение не должно заканчиваться на ! ? или ...
./examples/bad/main.go:27:2: лог-сообщение может содержать чувствительные данные: найдено ключевое слово "token"
```

### Хороший код (без нарушений)

```go
package main

import (
    "log/slog"
    "os"
    "go.uber.org/zap"
)

func main() {
    slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    zapLogger := zap.NewExample()

    userID := 12345

    slogLogger.Info("starting server on port 8080")
    slogLogger.Error("failed to connect to database")
    slogLogger.Debug("processing request", "user_id", userID)
    slogLogger.Warn("timeout occurred")

    zapLogger.Info("request completed")
    zapLogger.Error("connection refused")
    zapLogger.Debug("cache miss")
}
```

Вывод линтера: пусто (нет ошибок).

## 🧪 Запуск тестов

```bash
make test
# или
go test -v ./...
```

Тесты используют `analysistest` и проверяют все правила на примерах из `testdata/src/bad` и `testdata/src/good`.

## 📂 Структура проекта

```
.
├── cmd/
│   └── mylinter/           # точка входа
│       ├── main.go         # для go vet
│       └── plugin.go       # для golangci-lint (опционально)
├── pkg/
│   └── analyzer/           # основная логика
│       ├── analyzer.go
│       ├── analyzer_test.go
│       └── testdata/        # тестовые файлы
│           └── src/
│               ├── bad/     # плохие примеры
│               └── good/    # хорошие примеры
├── examples/                 # примеры использования
│   ├── bad/
│   │   └── main.go
│   └── good/
│       └── main.go
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## 📋 Команды Makefile

| Команда       | Описание                        |
|---------------|---------------------------------|
| `make build`  | сборка линтера (бинарник mylinter) |
| `make test`   | запуск тестов                   |
| `make vet`    | проверка примеров через go vet  |
| `make clean`  | удаление артефактов сборки      |
| `make all`    | полная сборка и проверка        |

## 📄 Лицензия

MIT License. Смотрите файл [LICENSE](LICENSE) для деталей.

---

**Примечание:** замените `yourusername` на актуальное имя пользователя или организации в ссылках на репозиторий.
```