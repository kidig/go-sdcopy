# SD Copy

SD Copy copies files from a source directory to a destination directory, preserving the directory structure and file timestamps.
It can be used for organizing destination files by date (see **Placeholders** below)

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/sdcopy.git
    cd sdcopy
    ```

2. Build the project:
    ```sh
    go build -o sdcopy main.go
    ```

## Usage

Run the `sdcopy` command with the source and destination paths as arguments:

```sh
./sdcopy /path/to/source /path/to/destination
```

This will copy all files from /path/to/source to /path/to/destination, preserving the directory structure and file timestamps.

### Placeholders in Destination Path
You can use placeholders in the destination path to dynamically create directories based on the file's modification time. The following placeholders are supported:

* {year}: The year of the file's modification time (e.g., 2023)
* {month}: The month of the file's modification time (e.g., 01)
* {day}: The day of the file's modification time (e.g., 15)

For example:
```sh
./sdcopy /path/to/source /path/to/destination/{year}/{month}/{day}
```

If a file was modified on Jan 15, 2023, it will be copied to:
```sh
/path/to/destination/2023/01/15/filename
```