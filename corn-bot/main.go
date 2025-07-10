package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var event_regex = regexp.MustCompile(`.*https://tickets.lakhta.events/event/([^"]+)".*`)
var interval = 5 * time.Second

type ScheduleResponse struct {
	Response ResponseResponse `json:"response"`
}

type ResponseResponse struct {
	Calendar []DayResponse `json:"calendar"`
}

type TimeSlot struct {
	Time     time.Time
	Quantity int
}

func (r *ResponseResponse) findTimeSlots() *[]TimeSlot {
	var timeSlots = []TimeSlot{}
	for _, dayResponse := range r.Calendar {
		slot_date, err := time.Parse("02.01.2006", dayResponse.Day)
		if err != nil {
			panic(err)
		}

		for _, timeResponse := range dayResponse.Time {
			quantity, err := strconv.Atoi(timeResponse.Quantity)
			if err != nil {
				panic(err)
			}

			if quantity > 0 {
				slot_time, err := time.Parse("15:04", timeResponse.Time)
				if err != nil {
					panic(err)
				}

				slot := time.Date(slot_date.Year(), slot_date.Month(), slot_date.Day(),
					slot_time.Hour(), slot_time.Minute(), 0, 0, time.Local)
				timeSlots = append(timeSlots, TimeSlot{Time: slot, Quantity: quantity})
			}

		}
	}
	return &timeSlots
}

type DayResponse struct {
	Day  string         `json:"day"`
	Time []TimeResponse `json:"_time"`
}

type TimeResponse struct {
	Time     string `json:"time"`
	Quantity string `json:"quantity"`
}

func main() {
	jar, err := cookiejar.New(nil) // Create a new in-memory cookie jar
	if err != nil {
		log.Panic().Err(err).Msg("Error creating cookie jar")
	}
	http.DefaultClient.Jar = jar

	for {
		event_id, err := getEventId()
		if err != nil {
			panic(err)
		}
		log.Info().Str("event_id", event_id).Msg("Got event id")

		schedule, err := getSchedule(event_id)
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Got schedule")

		slots := schedule.findTimeSlots()
		if len(*slots) > 0 {
			message := "Доступные слоты в Лахта Центр:\n"
			for _, slot := range *slots {
				message = message + slot.Time.Format("02 Jan, 15:04") + ": " + strconv.Itoa(slot.Quantity) + " мест\n"
			}
			err = sendTelegramNotification(message)
			if err != nil {
				panic(err)
			}
			err = sendTelegramNotification("https://tickets.lakhta.events/event/" + event_id)
			if err != nil {
				panic(err)
			}
			interval = 1 * time.Minute
		}

		time.Sleep(interval)
	}
}

func sendTelegramNotification(message string) error {
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("TELEGRAM_TOKEN")), nil)
	q := request.URL.Query()
	q.Add("chat_id", os.Getenv("TELEGRAM_CHAT_ID"))
	q.Add("text", message)
	request.URL.RawQuery = q.Encode()

	log.Info().Str("message", message).Msg("Sending telegram notification")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("Error calling telegram API")
	}
	response_body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading telegram response")
	}
	log.Info().Str("response", string(response_body)).Msg("Sent telegram notification")
	return err
}

func getEventId() (string, error) {
	main_page_reponse, err := http.DefaultClient.Get("https://lakhta.center/")
	if err != nil {
		log.Error().Err(err).Msg("Error calling main page https://lakhta.center/")
		return "", err
	}
	log.Info().Msg("Called main page https://lakhta.center/")

	main_page_content, err := io.ReadAll(main_page_reponse.Body)
	if err != nil {
		return "", err
	}

	var main_page_html = string(main_page_content)
	match := event_regex.FindStringSubmatch(main_page_html)
	return match[1], nil
}

func getSchedule(event_id string) (*ResponseResponse, error) {
	var body = fmt.Sprintf("{\"hash\":\"%s\"}", event_id)
	get, _ := http.NewRequest("GET", "https://tickets.lakhta.events/event/"+event_id, nil)
	get.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	response, err := http.DefaultClient.Do(get)
	if err != nil && response.StatusCode/100 != 2 {
		log.Error().Err(err).Msg("Error fetching main event page")
		return nil, err
	}

	post, _ := http.NewRequest("POST", "https://tickets.lakhta.events/api/no-scheme", strings.NewReader(body))
	post.Header.Add("Content-Type", "application/json")
	post.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")
	post.Header.Add("Referer", "https://tickets.lakhta.events/event/"+event_id)
	post.Header.Add("Origin", "https://tickets.lakhta.events")
	schedule_response, err := http.DefaultClient.Do(post)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching schedule")
		return nil, err
	}
	log.Info().Msg("Called schedule page https://tickets.lakhta.events/api/no-scheme")

	log.Info().Msg("Fetched schedule")

	schedule_response_content, err := io.ReadAll(schedule_response.Body)
	if err != nil {
		return nil, err
	}

	var scheduleResponse = ScheduleResponse{}
	err = json.Unmarshal(schedule_response_content, &scheduleResponse)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing schedule")
		return nil, err
	}

	return &scheduleResponse.Response, nil
}
