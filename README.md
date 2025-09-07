# rm-importer

https://github.com/user-attachments/assets/0996b7c3-85a7-45ea-90d7-cfabd989c5dc

A comprehensive tool for managing files on your reMarkable device with both export and import capabilities.

## Export Features
* Supports exporting as many folders & notes as you want;
* Can download both .pdf and .rmdoc;
* Retries the download **from the last failed note**;
* Waits for large notes long enough;
* Doesn't require reMarkable account or internet connection;
* Works with out of the box reMarkable software;
* Has a nice GUI.

## Import Features (NEW!, WINDOWS ONLY)
* **SSH-based file uploads** - Upload PDF files directly to your reMarkable device
* **Local file management** - Browse and select files from your computer
* **SSH file listing** - View all documents and folders on your device via SSH
* **Direct device integration** - Files are immediately available on your reMarkable
* **No cloud dependency** - All operations work locally without internet connection
* **Automatic metadata generation** - Creates proper metadata and content files for seamless integration

## Usage
Releases for Windows/MacOS/Linux are available on the 'Releases' tab of the repository.

The tool is built with [wailsv2](https://github.com/wailsapp/wails). The UI is implemented in Typescript/Svelte, file operations are done in Golang.

### Supported rM software version
Around 3.10+, around that version the local server requests got updated.

Tested on Version 3.16.2.3 on reMarkable 2.

### Prerequisites for Export (USB Mode)
* Enable USB connection in the Storage settings. Without the permission the app can't find the tablet;
* For long exports with large number of files, turn off Sleep Mode in the Battery settings. For some reason the local export doesn't prevent the tablet from going to sleep.

### Prerequisites for Import (SSH Mode)
* Enable SSH access on your reMarkable device (Developer Options → SSH)
* Install PuTTY's `plink` utility (included with PuTTY installation)
* Ensure your device and computer are on the same network
* Default SSH credentials: username `root`, password `[your device password]`

### Building steps
1. Install [wails v2](https://wails.io/docs/gettingstarted/installation).
2. Clone the project
3. `wails build`
