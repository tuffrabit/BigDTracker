# BigDTracker
CLI tool to keep track of your Devil's Ruin kills via the Bungie API

## Build
### Intro
Before you can use this, you must build it for yourself. This application makes use of the Bungie API which requires the use of a Bungie API key which is a private value. Once compiled, this application will contain a private Bungie API key embedded in the executable binary. You should not distribute an executable compiled with your private Bungie API key. Each person wanting to use this should register a new Bungie API app with Bungie to get their own private API key. Follow this guide to obtain your own private Bungie API key: https://github.com/vpzed/Destiny2-API-Info-wiki/blob/master/API-Introduction-Part-1-Setup.md#bungienet-api-key.

### Go
This application is written in Go. Building Go applications is easy:
1. Download and install Go for your OS: https://go.dev/dl/
2. Clone or download this repo (you might need to fork it for the Bungie API registration process)
3. Open a terminal and navigate to where you placed this repo
4. Run `go build -ldflags "-X main.apiKey=YOUR_BUNGIE_API_KEY_GOES_HERE"`

If everything went well you should have an executable with your private Bungie API embedded within.

## Usage
When running the application, you at least need to provide a Bungie name (username#5555). You may also provide an output location flag for debugging purposes.
### Windows (PowerShell)
No debug output `./BigDTracker username#5555`

With debug output `./BigDTracker username#5555 1`

Pipe debug output to file `./BigDTracker username#5555 1 | Out-File -FilePath .\out.txt`

## Notes
- As this application runs, it will store most data it gets from the Bungie API locally in a SQLite DB. That DB lives right next to the compiled executable. It does this because the Bungie API can be slow and we don't want to spam it. The first time you run this application, it will take some time (20 - 40 min depending on your PvP history). Subsequent runs will be much faster as it will use locally stored data when it can.