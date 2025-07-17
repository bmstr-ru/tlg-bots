package main

import (
	"bytes"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"strconv"
)

func NotifyTelegram(lot *Lot) error {
	b, err := bot.New(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		return err
	}

	chatId, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if lot.Image != nil {
		_, err = b.SendPhoto(context.TODO(), &bot.SendPhotoParams{
			ChatID:  chatId,
			Photo:   &models.InputFileUpload{Data: bytes.NewReader(lot.Image)},
			Caption: "https://torgi.gov.ru/new/public/lots/lot/" + lot.Id,
		})
	}

	return err
}
