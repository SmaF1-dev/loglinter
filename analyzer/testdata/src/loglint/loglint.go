package loglint

import (
	"log/slog"
)

const startMsg = "starting constant"
const badStartMsg = "Bad constant"

func test() {
	slog.Info("starting server")
	slog.Info("all good 123")

	slog.Info("Starting server") // want `log message should start with a lowercase letter`

	slog.Info("starting server!")   // want `log message should contain only English letters, digits, and spaces \(no special characters or emojis\)`
	slog.Info("starting server...") // want `log message should contain only English letters, digits, and spaces \(no special characters or emojis\)`
	slog.Info("starting сервер")    // want `log message should contain only English letters, digits, and spaces \(no special characters or emojis\)`

	slog.Info("user password: 123") // want `log message should contain only English letters, digits, and spaces \(no special characters or emojis\)` `log message may contain sensitive data: "password"`
	slog.Info("api_key = secret")   // want `log message should contain only English letters, digits, and spaces \(no special characters or emojis\)` `log message may contain sensitive data: "api_key"`

	slog.Info("token is abc") // want `log message may contain sensitive data: "token"`
	slog.Info("auth ok")      // want `log message may contain sensitive data: "auth"`

	slog.Info(startMsg)
	slog.Info(badStartMsg) // want `log message should start with a lowercase letter`
}
