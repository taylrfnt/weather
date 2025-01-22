# What Is This Project?
This is a simple little tool that I wrote to familiarize myself with how CLI applications
in Go work, as well as some GitHub Actions testing & publishing for Go.

This CLI tool provides you a binary (`weather`) that you can use to display the following:
- Active Weather Alerts
- Current Weather Information
- Forecasted Weather Conditions (hourly for the remainder of the day)

# Data Privacy / Disclosures
The `weather` binary will collect the following information:
- Your Public IP address
    - This is obtained using an API provided by [IPinfo](https://ipinfo.io).  This is a widely
    used public IP information service.
- An approximated location based on the Public IP Address obtained
    - The approximated location is stored as a set of coordinates (latitude & longitude), which
    are obtained by referencing the data from [ip-api](https://ip-api.com/) associated with the
    previously identified public IP address.

**This information is NOT stored anywhere but on your device by the `weather` binary.  We do not
control or manage the data policies of IPInfo or ip-api, and only use their services in alignment
with their non-commerical use agreements.**

# Limitations
- The source of the data for this CLI tool is the [US National Weather
Service's APIs](https://www.weather.gov/documentation/services-web-api).  As such, this
tool will only be able to provide weather information for locations within the United States.
- The data is currently **only** from the NWS's APIs, meaning that if the location does not have a
nearby station or data, this CLI tool will return no information.

# Installation
## Build from source
### Prerequisites
To build from source, you will need the following available on your host:
- `git`
- `go` (1.23+)
### Building
You can build this project from source by downloading a release tag artifact or cloning this repo:
```
git clone git@github.com:taylrfnt/weather.git weather
```
After cloning the repo, delete the git contents:
```
cd weather && rm -rf .git
```
Now, run a Go build:
```
go build
```
You can execute `weather` directly:
```
./weather
```
or you can add it to your `PATH` for easier use:
```
export PATH=${PATH}:/path/to/weather/directory
```

