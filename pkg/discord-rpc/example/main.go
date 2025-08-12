package main

import (
	"fmt"
	"time"

	"github.com/Voxelum/minecraft-launcher-core/pkg/discord-rpc/client"
)

func main() {
	err := client.Login("your-client-id")
	if err != nil {
		panic(err)
	}

	now := time.Now()
	timestamp := uint64(now.Unix())
	err = client.SetActivity(client.PayloadActivity{
		Type:    client.ActivityType.Playing,
		Details: "Playing Minecraft",
		State:   "In Game",
		Timestamps: &client.PayloadTimestamps{
			Start: &timestamp,
		},
	})

	if err != nil {
		panic(err)
	}

	// Discord will only show the presence if the app is running
	// Sleep for a few seconds to see the update
	fmt.Println("Sleeping...")
	time.Sleep(time.Second * 10)
}
