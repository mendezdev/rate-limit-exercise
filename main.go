package main

import (
	"fmt"

	"github.com/mendezdev/rate-limit-example/application"
)

/**
existing rate limit configurations:

-- STATUS 		-> 2 per 1 minute
-- NEWS 		-> 1 per 24 hours
-- MARKETING 	-> 3 per 1 hour

this can be found and changed in /mock_store/rate_limit_configuration_repo.json
*/

func main() {
	app := application.NewApp()

	// --------- sending STATUS notifications ---------
	app.SendNotification("1234", "STATUS", "STATUS message 1")
	app.SendNotification("1234", "STATUS", "STATUS message 2")
	// next one exceed rate limit config for STATUS
	app.SendNotification("1234", "STATUS", "STATUS message 3")

	// this HR notification does not have rate limit config so it will send it
	app.SendNotification("1234", "HR", "HR message 1")

	// --------- sending NEWS notifications ---------
	app.SendNotification("1234", "NEWS", "NEWS message 1")
	// next one exceed rate limit config for NEWS
	app.SendNotification("1234", "NEWS", "NEWS message 2")

	// --------- sending MARKETING notifications ---------
	app.SendNotification("1234", "MARKETING", "NEWS message 1")
	app.SendNotification("1234", "MARKETING", "NEWS message 2")
	app.SendNotification("1234", "MARKETING", "NEWS message 3")
	// next one exceed rate limit config for MARKETING
	app.SendNotification("1234", "MARKETING", "NEWS message 4")

	fmt.Println("program finished.")
}
