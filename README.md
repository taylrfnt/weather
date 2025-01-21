# Weather
## What Is This Project?
This is a simple little tool that I wrote to familiarize myself with how CLI applications
in Go work, as well as some GitHub Actions testing & publishing for Go.

This CLI tool provides you a binary (`weather`) that you can use to display the following:
- Active Weather Alerts
- Current Weather Information
- Forecasted Weather Conditions (hourly for the remainder of the day)

## Limitations
- The source of the data for this CLI tool is the [US National Weather
Service's APIs](https://www.weather.gov/documentation/services-web-api).  As such, this
tool will only be able to provide weather information for locations within the United States.
- The data is currently **only** from the NWS's APIs, meaning that if the location does not have a
nearby station or data, this CLI tool will return no information.

## TODO
1. Allow users provide a ZIP code for their location:
    - [ ] Option 1: Manually prompt
    - [ ] Option 2: pass as args
    - [ ] Option 3: config file(s)
2. Add tests for the following:
    - [ ] Incorrect input (invalid zip code)
    - [ ] Missing data
3. Figure out how to tag releases auto-magically with GHA

