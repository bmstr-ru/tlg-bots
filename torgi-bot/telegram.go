package main

import (
	"bytes"
	"context"
	"os"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func NotifyTelegram(lot *Lot) error {
	b, err := bot.New(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		return err
	}

	chatId, _ := strconv.ParseInt(os.Getenv("TG_CHANNEL_ID"), 10, 64)
	threadId, _ := strconv.ParseInt(os.Getenv("TG_THREAD_ID"), 10, 64)
	if lot.Image != nil {
		_, err = b.SendPhoto(context.TODO(), &bot.SendPhotoParams{
			ChatID:          chatId,
			MessageThreadID: int(threadId),
			Photo:           &models.InputFileUpload{Data: bytes.NewReader(lot.Image)},
			Caption:         "https://torgi.gov.ru/new/public/lots/lot/" + lot.Id,
		})
	}

	return err
}
