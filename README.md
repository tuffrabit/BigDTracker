# BigDTracker
CLI tool to keep track of your Devil's Ruin PvP kills via the Bungie API

## Bungie API Key
Before you can use this, you must aquire a Bungie API key for yourself and set it as an environment variable. Follow this guide to obtain your own private Bungie API key: https://github.com/vpzed/Destiny2-API-Info-wiki/blob/master/API-Introduction-Part-1-Setup.md#bungienet-api-key. Once you have a Bungie API key you must set a BUNGIE_API_KEY env variable:

### Windows (PowerShell)
`$env:BUNGIE_API_KEY="YOUR_BUNGIE_API_KEY_GOES_HERE"`

## Usage
When running the application, you at least need to provide a Bungie name (username#5555). You may also provide an output location flag for debugging purposes.

### Windows (PowerShell)
No debug output `./BigDTracker username#5555`

With debug output `./BigDTracker username#5555 1`

Pipe debug output to file `./BigDTracker username#5555 1 | Out-File -FilePath .\out.txt`

## Notes
- As this application runs, it will store most data it gets from the Bungie API locally in a SQLite DB. That DB lives right next to the compiled executable. It does this because the Bungie API can be slow and we don't want to spam it. The first time you run this application, it will take some time (20 - 40 min depending on your PvP history). Subsequent runs will be much faster as it will use locally stored data when it can.