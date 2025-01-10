# rm-exporter

https://github.com/user-attachments/assets/be7a0298-769c-4d6f-b366-477aef66df74

As you might know, reMarkable supports exporting notes locally through the USB connection.

Unfortunately, the default local export has a few flaws:
  * Large notes (10MB+) often can't be exported, **the UI doesn't wait for long enough for a note to download**;
  * Downloading a folder is not possible; only notes one-by-one.

This tool aims to fix those problems.

## Features
* supports exporting as many folders & notes as you want;
* can download both .pdf and .rmdoc;
* Retrying the download **from the last failed note**;
* Waits for large notes long enough;
* Doesn't require reMarkable account or internet connection;
* has a nice GUI.

## Usage
Releases for Windows/MacOS/Linux are available on the 'Releases' tab of the repository.

The tool is built with [wailsv2](https://github.com/wailsapp/wails). The UI is implemented in Typescript/Svelte, file downloading itself is done in Golang.

### Building steps
1. Install [wails v2](https://wails.io/docs/gettingstarted/installation).
2. Clone the project
3. `wails build`
