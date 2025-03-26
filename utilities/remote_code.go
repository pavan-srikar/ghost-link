package utilities

import (
	"bytes"
	"log"
	"os/exec"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ExecuteShellCommand runs a shell command and sends the output back to the user.
func ExecuteShellCommand(bot *tgbotapi.BotAPI, chatID int64, command string) {
	cmd := exec.Command("sh", "-c", command) // Use "cmd" for Windows
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Error executing command: "+stderr.String()))
		log.Println("Command execution error:", err)
		return
	}

	bot.Send(tgbotapi.NewMessage(chatID, "Command output:\n"+out.String()))
}

// ExecuteCode runs code in the specified language and sends the result back to the user.
func ExecuteCode(bot *tgbotapi.BotAPI, chatID int64, language, code string) {
	switch language {
	case "python":
		cmd := exec.Command("python3", "-c", code)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Python error: "+stderr.String()))
			log.Println("Python execution error:", err)
			return
		}

		bot.Send(tgbotapi.NewMessage(chatID, "Python output:\n"+out.String()))

	case "bash":
		cmd := exec.Command("bash", "-c", code)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Bash error: "+stderr.String()))
			log.Println("Bash execution error:", err)
			return
		}

		bot.Send(tgbotapi.NewMessage(chatID, "Bash output:\n"+out.String()))

	default:
		bot.Send(tgbotapi.NewMessage(chatID, "Unsupported language. Use 'python' or 'bash'."))
	}
}
