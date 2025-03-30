package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GeoResponse holds the geolocation response from the API
type GeoResponse struct {
	IP       string  `json:"ip"`
	City     string  `json:"city"`
	Region   string  `json:"region"`
	Country  string  `json:"country"`
	Loc      string  `json:"loc"`      // Format: "latitude,longitude"
	Org      string  `json:"org"`      // Internet Service Provider or organization
	Timezone string  `json:"timezone"` // Timezone
}

// HandleLocationCommand fetches and sends the device's location
func HandleLocationCommand(bot *tgbotapi.BotAPI, chatID int64) {
	apiURL := "https://ipinfo.io/json"

	// Make an HTTP GET request to fetch location details
	resp, err := http.Get(apiURL)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to get location: Unable to connect to the geolocation service."))
		return
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to get location: Unable to read the response."))
		return
	}

	var geoResponse GeoResponse
	if err := json.Unmarshal(body, &geoResponse); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Failed to get location: Unable to parse the response."))
		return
	}

	// Format the location details
	locationDetails := fmt.Sprintf("🌍 **Location Details**\n\n"+
		"**IP Address:** %s\n"+
		"**City:** %s\n"+
		"**Region:** %s\n"+
		"**Country:** %s\n"+
		"**Coordinates:** %s\n"+
		"**Organization:** %s\n"+
		"**Timezone:** %s\n"+
		"**Operating System:** %s",
		geoResponse.IP,
		geoResponse.City,
		geoResponse.Region,
		geoResponse.Country,
		geoResponse.Loc,
		geoResponse.Org,
		geoResponse.Timezone,
		runtime.GOOS)

	// Send the location details as a message
	msg := tgbotapi.NewMessage(chatID, locationDetails)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
