package main

import "time"

// Now that the service supports async calls, we should add a tmp value to the cache while the actual data is processed
// if the user tries to access the url before the resize job is done `tmpValue` will be displayed. We decide whether
// to write `terminalValue` (the actual resized picture) or `tmpValue` based on the boolean value `ongoingProcess`.
//
// In case the resize job for this entry fails, `terminalValue` will contain the error.
type cacheEntry struct {
	terminalValue []byte
	tmpValue []byte
	ongoingProcessing bool
	chMsg chan bool
	createdAt time.Time
}
