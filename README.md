# loglint

Линтер на Go для проверки лог‑сообщений. Совместим с golangci-lint (v2 custom binary) и поддерживает `log/slog` и `go.uber.org/zap` (Logger и SugaredLogger).

## Что проверяется

Проверки выполняются для **первого аргумента** вызова лог‑метода.

- `require-lowercase-start`: сообщение должно начинаться со строчной буквы (пробелы в начале игнорируются).
- `require-english`: допускаются только английские буквы (ASCII A–Z/a–z).
- `forbid-special-chars`: кроме букв/цифр/пробелов ничего не допускается (знаки пунктуации, `_`, `:`, `!`, эмодзи — запрещены).
- `forbid-sensitive-data`: запрещены чувствительные данные по ключевым словам и regex‑паттернам.

## Требования

- Go 1.24.0
- golangci-lint v2 (для режима custom binary)

## Быстрый старт (singlechecker)

1) Локальный запуск:
```bash
go run ./cmd/loglint ./examples/sample.go
```
2) Локальный запуск с исправлениями:
```bash
go run ./cmd/loglint -diff -fix ./examples/sample.go
```


Примечание: singlechecker запускается **с дефолтными настройками** (все проверки включены). CLI‑флагов для настройки нет — конфигурация задаётся через golangci-lint.

## Интеграция с golangci-lint (v2 custom binary)

1) Соберите кастомный бинарник:

```bash
golangci-lint custom
```

Команда читает `.custom-gcl.yml` (указаны модуль и импорт):

```yaml
version: v2.9.0
plugins:
  - module: example.com/loglint
    import: example.com/loglint/pkg/loglint
    path: .
```

2) Настройте линтер в `.golangci.yml`:

```yaml
version: "2"

linters:
  default: none
  enable:
    - loglint

  settings:
    custom:
      loglint:
        type: module
        description: Checks log message rules (lowercase, English, special chars, sensitive data)
        settings:
          require-lowercase-start: true
          require-english: true
          forbid-special-chars: true
          forbid-sensitive-data: true
          sensitive-keywords:
            - password
            - api_key
            - token
          sensitive-patterns:
            - (?i)secret\\s*[:=]
```

3) Запуск:

```bash
./custom-gcl run ./examples/sample.go
```

## Конфигурация

Все опции задаются в `linters.settings.custom.loglint.settings`:

- `require-lowercase-start` (bool)
- `require-english` (bool)
- `forbid-special-chars` (bool)
- `forbid-sensitive-data` (bool)
- `sensitive-keywords` ([]string)
- `sensitive-patterns` ([]string, regex)

## Примеры

Смотри `examples/sample.go`.

## Тесты

```bash
go test ./...
```

Тесты используют `analysistest`, данные в `pkg/loglint/testdata`.

## Соответствие заданию

- Реализация правил (строчная буква, английский, спецсимволы, чувствительные данные): готово.
- Тесты для правил: готово.
- Интеграция с golangci-lint (custom binary): готово.
- Бонусы:
  - Конфигурация правил: готово.
  - SuggestedFixes: реализован авто‑фикс для `require-lowercase-start` (для строковых литералов).
  - Кастомные паттерны для чувствительных данных (regex): готово.
  - CI: готово.

## Ограничения

- Анализируется **только первый аргумент** лог‑вызова.
- Для `zap.SugaredLogger` проверяются только методы `Infof/Debugf/...` и `Infow/Debugw/...` (у них первый аргумент - сообщение). Методы `Info/Debug/...` игнорируются.
- Правило `forbid-special-chars` строгое (запрещает любую пунктуацию).

## Как помог ChatGPT

1. Подсказал варианты реализации правил и подсветил edge‑кейсы для тестов. Реализовал сами тесты
2. Обсудили конфигурацию линтера и подход к проверке чувствительных данных (ключевые слова + regex). Реализовал проверку через regex.
3. Провел финальную проверку тестов на необработанные кейсы и соответствие изначальному заданию
4. Проверил реализацию SuggestedFixes и предложил улучшения с пояснениями + помог реализовать тесты для этого случая
5. В процессе доработки правил подсказал конкретные стандартные методы/подходы в Go, которые подходят для этих проверок