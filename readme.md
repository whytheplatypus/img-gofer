# img-gofer

`img-gofer` is a command-line interface (CLI) tool written in Go that allows you to download your Google Photos library to the current directory for backup purposes. With `img-gofer`, you can easily preserve your photos and videos from Google Photos on a local device, ensuring that you have a backup in case of data loss or other issues.

## Installation

To install `img-gofer`, simply download the binary file for your platform from [the latest release](https://github.com/whytheplatypus/img-gofer/releases), and add it to your PATH. 

## Usage

Before running `img-gofer`, you must set the `CLIENT_ID` and `CLIENT_SECRET` environment variables to your [Google API credentials](https://developers.google.com/photos/library/guides/get-started#enable-the-api) for your Google Photos account.

`img-gofer` can be run with the following command:

```bash
img-gofer
```

`img-gofer` will download all of the photos in your Google Photos library to the current directory. It will skip any files that are already present in the current directory to avoid duplicating files.

## Contributing

If you'd like to contribute to `img-gofer`, please fork the repository and submit a pull request.