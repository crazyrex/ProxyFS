package main

func fetchAuthTokenAndAccountURL(forceUpdate bool) (swiftAuthToken string, swiftAccountURL string) {
	// Handle typical case where both are already available (unless update is forced)

	globals.Lock()
	if !forceUpdate && ("" != globals.swiftAuthToken) && ("" != globals.swiftAccountURL) {
		swiftAuthToken = globals.swiftAuthToken
		swiftAccountURL = globals.swiftAccountURL
		globals.Unlock()
		return
	}
	globals.Unlock()

	// TODO
	swiftAuthToken = ""
	swiftAccountURL = ""
	return
}
