package main

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	lots, err := FetchLots()

	if err != nil {
		log.Panic().Err(err).Msg("Error fetching torgi lots")
		os.Exit(1)
	}

	for _, lot := range lots {
		log.Info().Any("lot", lot).Msg("Fetched lot")
		if LotChanged(lot) {
			err = NotifyTelegram(lot)
			if err != nil {
				log.Panic().Err(err).Msg("Error notifying on lot change")
				os.Exit(1)
			}
			err = StoreLot(lot)
			if err != nil {
				log.Panic().Err(err).Msg("Error storing information about lot")
				os.Exit(1)
			}
		}
	}

}

func FetchLots() ([]*Lot, error) {
	lotsRequest, _ := http.NewRequest("GET", "https://torgi.gov.ru/new/api/public/lotcards/search", nil)
	q := lotsRequest.URL.Query()
	q.Add("catCode", "2")
	q.Add("fiasGUID", "e6010c50-dfbb-4395-b68b-ed03bedc0c1e")
	q.Add("byFirstVersion", "true")
	q.Add("withFacets", "true")
	q.Add("page", "0")
	q.Add("size", "10")
	q.Add("sort", "firstVersionPublicationDate,desc")
	lotsRequest.URL.RawQuery = q.Encode()

	lotsResponse, err := http.DefaultClient.Do(lotsRequest)
	if err != nil {
		return nil, err
	}
	defer lotsResponse.Body.Close()

	lotsResponseBody, err := io.ReadAll(lotsResponse.Body)
	if err != nil {
		return nil, err
	}
	log.Info().Str("body", string(lotsResponseBody)).Msg("Successfully requested lots")

	lotsResponseStruct := ApiResponse{}
	err = json.Unmarshal(lotsResponseBody, &lotsResponseStruct)
	if err != nil {
		return nil, err
	}

	log.Info().Int("count", len(lotsResponseStruct.Content)).Msg("Parsed response body")

	var lots []*Lot

	for _, apiLot := range lotsResponseStruct.Content {
		var image []byte
		if len(apiLot.LotImages) > 0 {
			image, err = FetchImage(apiLot.LotImages[0])
			if err != nil {
				log.Error().Err(err).Msg("Could not fetch lot image")
			}
		}
		createDate, err := time.Parse(time.RFC3339Nano, apiLot.CreateDate)
		if err != nil {
			log.Error().Err(err).Str("createDate", apiLot.CreateDate).Msg("Could not parse createDate")
		}

		bidEndTime, err := time.Parse(time.RFC3339Nano, apiLot.BiddEndTime)
		if err != nil {
			log.Error().Err(err).Str("biddEndTime", apiLot.BiddEndTime).Msg("Could not parse biddEndTime")
		}

		lot := &Lot{
			Id:          apiLot.ID,
			Status:      apiLot.LotStatus,
			Name:        apiLot.LotName,
			Description: apiLot.LotDescription,
			Image:       image,
			CreateDate:  createDate,
			BidEndTime:  bidEndTime,
			HasAppeals:  apiLot.HasAppeals,
			IsAnnulled:  apiLot.IsAnnulled,
		}
		lots = append(lots, lot)
	}

	return lots, nil
}

func FetchImage(imageId string) ([]byte, error) {
	response, err := http.DefaultClient.Get("https://torgi.gov.ru/new/image-preview/v1/" + imageId)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}
