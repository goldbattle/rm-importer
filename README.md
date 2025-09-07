# rm-importer

<img width="50%" alt="example_00" src="https://github.com/user-attachments/assets/89d08f9f-2bbe-4311-baac-2bd83098bb42" /><img width="50%" alt="example_01" src="https://github.com/user-attachments/assets/4a9d18c0-1727-49bd-b941-4c93ea31e024" />


A comprehensive tool for managing files on your reMarkable device with both export and import capabilities.

## Import Features (NEW!, WINDOWS ONLY)
* **100% vibe coded** - Don't count on much, but it seems to work?
* **SSH-based file uploads** - Upload PDF files directly to your reMarkable device
* **SSH file listing** - View all documents and folders on your device via SSH instead of HTTP
* **Direct device integration** - Files are immediately available on your reMarkable
* **Automatic metadata generation** - Creates proper metadata and content files so they show up in the on-device GUI


## Export Features
* Supports exporting as many folders & notes as you want;
* Can download both .pdf and .rmdoc;
* Retries the download **from the last failed note**;
* Waits for large notes long enough;
* Doesn't require reMarkable account or internet connection;
* Works with out of the box reMarkable software;
* Has a nice GUI.


## Usage
Releases for Windows/MacOS/Linux are available on the 'Releases' tab of the repository.

The tool is built with [wailsv2](https://github.com/wailsapp/wails). The UI is implemented in Typescript/Svelte, file operations are done in Golang.

### Supported rM software version
Tested on Version 3.3.2.1666 on reMarkable 2.

### Prerequisites for Import (SSH Mode)
* Ensure your device and computer are on the same network or use USB cable
* Default SSH credentials: username `root`, password `[your device password]`
  * NOTE: You can find your IP and ssh pass in `Help > Copyrights and licenses" at the very bottom
  * https://remarkable.guide/guide/access/ssh.html
* Install PuTTY's `plink` utility (included with PuTTY installation)

### Prerequisites for Export (USB Mode)
* Enable USB connection in the Storage settings. Without the permission the app can't find the tablet;
* For long exports with large number of files, turn off Sleep Mode in the Battery settings. For some reason the local export doesn't prevent the tablet from going to sleep.


### Building steps
1. Install [wails v2](https://wails.io/docs/gettingstarted/installation).
2. Clone the project
3. `wails build`
4. Launch via `./build/bin/rm-importer.exe` after build
5. Debug via `wails dev` to get logs



