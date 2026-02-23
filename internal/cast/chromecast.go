package cast

import (
	"context"
	"fmt"
	"time"

	"github.com/vishen/go-chromecast/application"
)

// Play streams audio to a Chromecast device
func Play(ip string, port int, audioURL string, block bool) error {
	// Create application
	app := application.NewApplication()

	// Connect to device
	if err := app.Start(ip, port); err != nil {
		return fmt.Errorf("failed to connect to device: %w", err)
	}
	defer app.Close(true) // stopMedia = true

	// Load media: filenameOrUrl, startTime, contentType, transcode, detach, forceDetach
	if err := app.Load(audioURL, 0, "audio/mpeg", false, false, false); err != nil {
		return fmt.Errorf("failed to load media: %w", err)
	}

	if !block {
		return nil
	}

	// Wait for playback to complete
	return waitForPlayback(app)
}

func waitForPlayback(app *application.Application) error {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Give it a moment to start playing
	time.Sleep(500 * time.Millisecond)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("playback timeout")
		case <-ticker.C:
			// Status returns: *cast.Application, *cast.Media, *cast.Volume
			_, media, _ := app.Status()

			// If no media status, playback is done
			if media == nil {
				return nil
			}

			// Check player state
			playerState := media.PlayerState
			if playerState == "IDLE" || playerState == "" {
				return nil
			}
		}
	}
}
