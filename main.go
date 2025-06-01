package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	dotenv "github.com/joho/godotenv"

	"templates/modules"
)

var (
	devID    int64
	appID    int
	appHash  string
	botToken string
)

func init() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	devIDStr := os.Getenv("DEV_ID")
	parsedDevID, err := strconv.ParseInt(devIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid DEV_ID: %v", err)
	}
	devID = parsedDevID

	appIDStr := os.Getenv("APP_ID")
	appID, err = strconv.Atoi(appIDStr)
	if err != nil {
		log.Fatalf("Invalid APP_ID: %v", err)
	}

	appHash = os.Getenv("APP_HASH")
	botToken = os.Getenv("BOT_TOKEN")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	sessionDir := ".td"
	err := os.MkdirAll(sessionDir, 0700)
	if err != nil {
		log.Fatalf("Failed to create session dir: %v", err)
	}

	sessionStorage := &telegram.FileSessionStorage{
		Path: filepath.Join(sessionDir, "session.json"),
	}

	dispatcher := tg.NewUpdateDispatcher()
	client := telegram.NewClient(appID, appHash, telegram.Options{
		UpdateHandler:  dispatcher,
		SessionStorage: sessionStorage,
	})

	stop, err := bg.Connect(client)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = stop()
	}()

	if _, err := client.Auth().Bot(ctx, botToken); err != nil {
		log.Fatal(err)
	}

	me, err := client.Self(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Logged in as @%s", me.Username)

	api := client.API()
	modules.InitModules(ctx, api, dispatcher)

	log.Println("✅ Bot started. Press Ctrl+C to stop...")
	<-ctx.Done()
	log.Println("🛑 Interrupt signal received. Shutting down gracefully...")
}
