package hue

import (
	"errors"
	"fmt"
	"time"
)

const (
	// NoUpdateAvailable represents an up-to-date bridge
	NoUpdateAvailable       = 0
	// DownloadingSystemUpdate represents a bridge downloading a system update
	DownloadingSystemUpdate = 1
	// SystemUpdateAvailable represents a bridge with a system update ready to be applied
	SystemUpdateAvailable   = 2
	// SystemUpdating represents a bridge that is actively updating itself
	SystemUpdating          = 3
	// LastRequestFailed represents a bridge that failed to execute its previous command
	LastRequestFailed       = 10
	// NetworkUnavailable represents a bridge unable to update because it lacks Internet connectivity
	NetworkUnavailable      = 11
)

// Updater is an instance of a Hue bridge updater.
type Updater struct {
	bridge *Bridge
	state  int32

	checkInProgress bool

	err error
	msg string
}

// NewUpdater creates a new bridge updater from the specified Hue bridge.
func NewUpdater(b *Bridge) Updater {
	return Updater{
		bridge: b,
		state:  NoUpdateAvailable,
	}
}

// State exposes the current state of the updater.
func (u *Updater) State() int32 {
	return u.state
}

// Run begins the process of monitoring a bridge for updates then applying them.
func (u *Updater) Run(results chan string, quit chan interface{}) {
	ticker := time.NewTicker(60 * time.Minute)

	for {
		select {
		case <-ticker.C:
			switch u.state {
			case NoUpdateAvailable, DownloadingSystemUpdate:
				newState := u.checkForUpdate()

				if newState == u.state {
					break
				} else if newState == LastRequestFailed {
					results <- "Error checking for updates: " + u.err.Error()
					break
				} else if newState == DownloadingSystemUpdate {
					results <- u.msg
				} else if newState == SystemUpdateAvailable {
					results <- u.msg

					newState = u.executeUpdate()

					if newState == LastRequestFailed {
						results <- "Error applying update: " + u.err.Error()
						u.state = SystemUpdateAvailable
						break
					} else {
						results <- u.msg
					}
				}

				u.state = newState

			case SystemUpdateAvailable:
				newState := u.executeUpdate()

				if newState == LastRequestFailed {
					results <- "Error applying update: " + u.err.Error()
				} else {
					results <- u.msg
					u.state = newState
				}

			case SystemUpdating:
			default:
			}

		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (u *Updater) checkForUpdate() int32 {
	config, err := u.bridge.Config()
	if err != nil {
		u.err = err
		return LastRequestFailed
	} else if !config.PortalState.SignedOn {
		return NetworkUnavailable
	}

	if config.SwUpdate.State == 0 || config.SwUpdate.State == 1 {
		err = u.bridge.CheckForUpdate()
		if err != nil {
			u.err = err
			return LastRequestFailed
		}

		itr := 0
		for {
			time.Sleep(time.Minute)

			checkConfig, err := u.bridge.Config()

			if err != nil {
				u.err = err
				return LastRequestFailed
			}

			if checkConfig.SwUpdate.CheckForUpdates && itr < 5 {
				itr++
				continue
			} else if checkConfig.SwUpdate.CheckForUpdates {
				return NoUpdateAvailable
			}

			// Here checkForUpdates has gone back to false, and we aren't past the timeout,
			// which means the process can continue.
			config = checkConfig
			break
		}
	}

	switch config.SwUpdate.State {
	case 0:
		u.msg = ""
		return NoUpdateAvailable
	case 1:
		u.msg = "Downloading update"
		return DownloadingSystemUpdate
	case 2:
		u.msg = "Update available: " + config.SwUpdate.UpdateSummary
		return SystemUpdateAvailable
	case 3:
		u.msg = "Applying " + config.SwUpdate.UpdateSummary
		return SystemUpdating
	default:
		u.msg = fmt.Sprintf("Unknown state detected: %d\n", config.SwUpdate.State)
		return LastRequestFailed
	}
}

func (u *Updater) executeUpdate() int32 {
	startingConfig, err := u.bridge.Config()
	if err != nil {
		u.err = err
		return LastRequestFailed
	}

	u.bridge.updateInProgress = true
	defer func() {
		u.bridge.updateInProgress = false
	}()

	err = u.bridge.StartUpdate()
	if err != nil {
		u.err = err
		return LastRequestFailed
	}

	u.state = SystemUpdating

	itr := 0
	for {
		time.Sleep(time.Minute)

		checkConfig, err := u.bridge.Config()

		// TODO: further error checking to differentiate between 'bridge is unavailable because rebooting' and 'bridge having true issues'
		if err != nil && itr < 5 {
			itr++
			continue
		} else if err != nil {
			u.err = err
			return LastRequestFailed
		}

		if !checkConfig.SwUpdate.NotifyUser && itr < 5 {
			itr++
			continue
		} else if !checkConfig.SwUpdate.NotifyUser {
			// We assume that things have changed, even if we don't explicitly detect it, when the API versions change.
			if checkConfig.SwVersion == startingConfig.SwVersion {
				u.err = errors.New("Failed to update")
				return LastRequestFailed
			}
		}

		// We either have detected that we need to notify the user, or we have detected a new API version.
		// Reset the notify flag then return.

		// We don't care if it fails; we still consider the update to be complete.
		_ = u.bridge.FinishUpdate()

		u.msg = "Bridge updated to " + checkConfig.APIVersion
		return checkConfig.SwUpdate.State
	}

	u.err = errors.New("unknown state detected")
	return LastRequestFailed
}
