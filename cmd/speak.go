package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Ozdotdotdot/cast/internal/cast"
	"github.com/Ozdotdotdot/cast/internal/config"
	"github.com/Ozdotdotdot/cast/internal/server"
	"github.com/Ozdotdotdot/cast/internal/tts"
)

func speak(cfg *config.Config, deviceSpec, message string) error {
	// Resolve device(s) from spec (could be device name, alias, or group)
	devices, err := resolveDevices(cfg, deviceSpec)
	if err != nil {
		return err
	}

	// Generate TTS audio
	audioPath, err := tts.Generate(message)
	if err != nil {
		return fmt.Errorf("failed to generate TTS: %w", err)
	}
	defer tts.Cleanup(audioPath)

	// Start local HTTP server
	srv, url, err := server.Start(audioPath)
	if err != nil {
		return fmt.Errorf("failed to start audio server: %w", err)
	}
	defer srv.Shutdown()

	// Cast to all devices
	if len(devices) == 1 {
		return castToDevice(devices[0], url, cfg.BlockUntilComplete)
	}

	// Multiple devices - cast in parallel
	var wg sync.WaitGroup
	errChan := make(chan error, len(devices))

	for _, d := range devices {
		wg.Add(1)
		go func(device config.Device) {
			defer wg.Done()
			if err := castToDevice(device, url, cfg.BlockUntilComplete); err != nil {
				errChan <- fmt.Errorf("%s: %w", device.OriginalName, err)
			}
		}(d)
	}

	wg.Wait()
	close(errChan)

	// Collect errors
	var errs []string
	for err := range errChan {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to cast to some devices:\n  %s", strings.Join(errs, "\n  "))
	}

	return nil
}

func castToDevice(device config.Device, audioURL string, block bool) error {
	return cast.Play(device.IP, device.Port, audioURL, block)
}

func resolveDevices(cfg *config.Config, spec string) ([]config.Device, error) {
	var devices []config.Device

	// Check if it's a comma-separated list
	parts := strings.Split(spec, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Check if it's a group
		if group, ok := cfg.Groups[part]; ok {
			for _, deviceName := range group.Devices {
				if d, err := cfg.GetDevice(deviceName); err == nil {
					devices = append(devices, d)
				}
			}
			continue
		}

		// Try to find as device or alias
		d, err := cfg.GetDevice(part)
		if err != nil {
			return nil, fmt.Errorf("unknown device or group: %s", part)
		}
		devices = append(devices, d)
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices resolved from: %s", spec)
	}

	return devices, nil
}
