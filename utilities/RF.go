package utilities

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var currentDir = "."

func HandleListCommand(bot *tgbotapi.BotAPI, chatID int64) {
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to list directory: %v", err)))
		return
	}

	var response string
	for _, entry := range entries {
		if entry.IsDir() {
			response += "[D] " + entry.Name() + "\n"
		} else {
			response += "[F] " + entry.Name() + "\n"
		}
	}

	if response == "" {
		response = "No files or folders found."
	}

	bot.Send(tgbotapi.NewMessage(chatID, response))
}

func HandleMoveToFolderCommand(bot *tgbotapi.BotAPI, chatID int64, folder string) {
	newPath := filepath.Join(currentDir, folder)
	if err := os.Chdir(newPath); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to move to folder: %v", err)))
		return
	}
	currentDir, _ = os.Getwd()
	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Moved to folder: %s", currentDir)))
}

func HandleMoveBackCommand(bot *tgbotapi.BotAPI, chatID int64) {
	if err := os.Chdir(".."); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to move back: %v", err)))
		return
	}
	currentDir, _ = os.Getwd()
	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Moved to folder: %s", currentDir)))
}

func HandleGetCommand(bot *tgbotapi.BotAPI, chatID int64, name string) {
	path := filepath.Join(currentDir, name)
	info, err := os.Stat(path)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to access file/folder: %v", err)))
		return
	}

	if info.IsDir() {
		sendFilesRecursively(bot, chatID, path)
	} else {
		sendFile(bot, chatID, path)
	}
}

func sendFile(bot *tgbotapi.BotAPI, chatID int64, filePath string) {
	file := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(filePath))
	_, err := bot.Send(file)
	if err != nil {
		log.Printf("Failed to send file: %s, error: %v", filePath, err)
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to send file: %s", filePath)))
	}
}

func sendFilesRecursively(bot *tgbotapi.BotAPI, chatID int64, folderPath string) {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking file: %v", err)
			return nil
		}

		if !info.IsDir() {
			sendFile(bot, chatID, path)
		}
		return nil
	})

	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to send folder: %v", err)))
	}
}
