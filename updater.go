package hue_go

import (
	"errors"
	"fmt"
	"time"
)

const (
	NO_UPDATE_AVAILABLE       = 0
	DOWNLOADING_SYSTEM_UPDATE = 1
	SYSTEM_UPDATE_AVAILABLE   = 2
	SYSTEM_UPDATING           = 3
	LAST_REQUEST_FAILED       = 10
	NETWORK_UNAVAILABLE       = 11
)

type Updater struct {
	bridge *Bridge
	state  int32

	checkInProgress bool

	err error
	msg string
}

func NewUpdater(b *Bridge) Updater {
	return Updater{
		bridge: b,
		state:  NO_UPDATE_AVAILABLE,
	}
}

func (u *Updater) State() int32 {
	return u.state
}

func (u *Updater) Run(results chan string, quit chan interface{}) {
	ticker := time.NewTicker(60 * time.Minute)

	for {
		select {
		case <-ticker.C:
			switch u.state {
			case NO_UPDATE_AVAILABLE, DOWNLOADING_SYSTEM_UPDATE:
				newState := u.checkForUpdate()

				if newState == u.state {
					break
				} else if newState == LAST_REQUEST_FAILED {
					results <- "Error checking for updates: " + u.err.Error()
					break
				} else if newState == DOWNLOADING_SYSTEM_UPDATE {
					results <- u.msg
				} else if newState == SYSTEM_UPDATE_AVAILABLE {
					results <- u.msg

					newState = u.executeUpdate()

					if newState == LAST_REQUEST_FAILED {
						results <- "Error applying update: " + u.err.Error()
						u.state = SYSTEM_UPDATE_AVAILABLE
						break
					} else {
						results <- u.msg
					}
				}

				u.state = newState

			case SYSTEM_UPDATE_AVAILABLE:
				newState := u.executeUpdate()

				if newState == LAST_REQUEST_FAILED {
					results <- "Error applying update: " + u.err.Error()
				} else {
					results <- u.msg
					u.state = newState
				}

			case SYSTEM_UPDATING:
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
		return LAST_REQUEST_FAILED
	} else if !config.PortalState.SignedOn {
		return NETWORK_UNAVAILABLE
	}

	if config.SwUpdate.State == 0 || config.SwUpdate.State == 1 {
		err = u.bridge.CheckForUpdate()

		if err != nil {
			u.err = err
			return LAST_REQUEST_FAILED
		}

		itr := 0
		for {
			time.Sleep(time.Minute)

			checkConfig, err := u.bridge.Config()

			if err != nil {
				u.err = err
				return LAST_REQUEST_FAILED
			}

			if checkConfig.SwUpdate.CheckForUpdates && itr < 5 {
				itr++
				continue
			} else if checkConfig.SwUpdate.CheckForUpdates {
				return NO_UPDATE_AVAILABLE
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
		return NO_UPDATE_AVAILABLE
	case 1:
		u.msg = "Downloading update"
		return DOWNLOADING_SYSTEM_UPDATE
	case 2:
		u.msg = "Update available: " + config.SwUpdate.UpdateSummary
		return SYSTEM_UPDATE_AVAILABLE
	case 3:
		u.msg = "Applying " + config.SwUpdate.UpdateSummary
		return SYSTEM_UPDATING
	default:
		u.msg = fmt.Sprintf("Unknown state detected: %d\n", config.SwUpdate.State)
		return LAST_REQUEST_FAILED
	}
}

func (u *Updater) executeUpdate() int32 {
	startingConfig, err := u.bridge.Config()

	if err != nil {
		u.err = err
		return LAST_REQUEST_FAILED
	}

	u.bridge.updateInProgress = true
	defer func() {
		u.bridge.updateInProgress = false
	}()

	err = u.bridge.StartUpdate()

	if err != nil {
		u.err = err
		return LAST_REQUEST_FAILED
	}

	u.state = SYSTEM_UPDATING

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
			return LAST_REQUEST_FAILED
		}

		if !checkConfig.SwUpdate.NotifyUser && itr < 5 {
			itr++
			continue
		} else if !checkConfig.SwUpdate.NotifyUser {
			// We assume that things have changed, even if we don't explicitly detect it, when the API versions change.
			if checkConfig.SwVersion == startingConfig.SwVersion {
				u.err = errors.New("Failed to update")
				return LAST_REQUEST_FAILED
			}
		}

		// We either have detected that we need to notify the user, or we have detected a new API version.
		// Reset the notify flag then return.

		// We don't care if it fails; we still consider the update to be complete.
		_ = u.bridge.FinishUpdate()

		u.msg = "Bridge updated to " + checkConfig.ApiVersion
		return checkConfig.SwUpdate.State
	}

	u.err = errors.New("Unknown state detected")
	return LAST_REQUEST_FAILED
}
