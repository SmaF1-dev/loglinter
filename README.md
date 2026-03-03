# loglinter

Линтер для Go, проверяющий лог-записи в коде на соответствие правилам оформления и безопасности.  
Поддерживаемые логгеры: 
- `log/slog`
- `go.uber.org/zap`

## Правила

1. Сообщение должно начинаться со **строчной буквы**.
2. Сообщение должно содержать только **английские буквы, цифры и пробелы**.
3. Сообщение не должно содержать **потенциально чувствительные данные** (по умолчанию: `password`, `api_key`, `token`, `secret`, `key`, `auth`).

## Возможности
- **Интеграция с golangci-lint** - может работать как плагин.
- **Автоисправление** - для правила о первой букве (флаг `-fix`).
- **Конфигурация** - настройка чувствительных слов через конфиг.

## Установка

### Самостоятельная утилита

```bash
go install github.com/SmaF1-dev/loglinter/cmd/loglinter@latest
```

### Плагин для golangci-lint

Клонируйте репозиторий и соберите плагин:

```bash
git clone https://github.com/SmaF1-dev/loglinter
cd loglinter
go build -buildmode=plugin -o loglinter.so ./cmd/plugin
```

Поместите loglinter.so в удобное место и в файле `.golangci.yml` укажите путь к плагину. 
При необходимости, можете указать настройки:
```
yaml
linters-settings:
  custom:
    loglinter:
      path: /path/to/loglinter.so
      settings:
        sensitive_keywords:
          - "password"
          - "token"
          - "api_key"
linters:
  enable:
    - loglinter
```

## Использование

### Самостоятельный запуск

```bash
loglinter ./...
```

Для автоматического исправления ошибок с первой буквой добавьте флаг -fix:

```bash
loglinter -fix ./...
```

### В составе golangci-lint

```bash
golangci-lint run
```

## Пример

```golang
package main

import "log/slog"

func main() {
    slog.Info("starting server") // OK
    slog.Info("Starting server") // ошибка: первая буква заглавная
    slog.Info("starting server!") // ошибка: спецсимвол
    slog.Info("user password: secret") // ошибка: чувствительные данные
}
```