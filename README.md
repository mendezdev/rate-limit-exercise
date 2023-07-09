# rate-limit-exercise
This is a simple example of notification system with rate limit

## Setup

You need golang version 1.20+

Here you can download [golang](https://go.dev/doc/install).

## Tests
There are a couple of test in the notification_service just to cover some flows.

Run ``` go test ./... ``` on terminal standing in the root folder of the app to execute all the tests.

## How to run?
Open the terminal, go to the app root folder and execute ``` go run main.go ```

## NOTES
There is a simple rate-limit setting that the app will pull up on startup and can be found at ```./mock_store/rate_limit_configuration.json```

I took this exercise as a simple first iteration to show a bit of code.

It is for this reason that some possible errors are simply logged because I imagine a scenario in which we agree that the flow of sending a notification should not be interrupted by failing to find a configuration among other things.

There are several improvement opportunities such as the ```user_id``` could be a unique identifier like a ```UUID```. There are no certain validations on the values ​​that can arrive in the ```NotificationRequest``` among other things that we could talk about in a code review meet.
