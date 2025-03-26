package utilities

import (
	"bytes"
	"image/png"
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kbinani/screenshot"
)

var screenshotFrequency time.Duration
var autoScreenshotEnabled = false

// HandleScreenshot captures a screenshot and sends it to the user.
func HandleScreenshot(bot *tgbotapi.BotAPI, chatID int64) {
	numDisplays := screenshot.NumActiveDisplays()
	if numDisplays <= 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "No active displays found."))
		return
	}

	// Capture the first display
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to capture screenshot."))
		log.Println("Error capturing screenshot:", err)
		return
	}

	// Encode screenshot as PNG
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to encode screenshot."))
		log.Println("Error encoding screenshot:", err)
		return
	}

	// Send screenshot as a photo
	file := tgbotapi.FileBytes{Name: "screenshot.png", Bytes: buf.Bytes()}
	bot.Send(tgbotapi.NewPhoto(chatID, file))
}

// SetScreenshotFrequency sets the frequency for auto-screenshots.
func SetScreenshotFrequency(bot *tgbotapi.BotAPI, chatID int64, freq string) {
	duration, err := time.ParseDuration(freq + "s")
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Invalid frequency format. Use seconds, e.g., /screenshot_frequency 30"))
		return
	}
	screenshotFrequency = duration
	bot.Send(tgbotapi.NewMessage(chatID, "Screenshot frequency set to "+freq+" seconds."))
}

// StartAutoScreenshot starts the auto-screenshot process.
func StartAutoScreenshot(bot *tgbotapi.BotAPI, chatID int64) {
	if screenshotFrequency == 0 {
		screenshotFrequency = 30 * time.Second // Default to 30 seconds
	}

	autoScreenshotEnabled = true
	bot.Send(tgbotapi.NewMessage(chatID, "Auto-screenshot started."))

	go func() {
		for autoScreenshotEnabled {
			HandleScreenshot(bot, chatID)
			time.Sleep(screenshotFrequency)
		}
	}()
}

// StopAutoScreenshot stops the auto-screenshot process.
func StopAutoScreenshot(bot *tgbotapi.BotAPI, chatID int64) {
	autoScreenshotEnabled = false
	bot.Send(tgbotapi.NewMessage(chatID, "Auto-screenshot stopped."))
}
