# varip

A lightweight command line environment variables ripper (think simple Grep for config files)

Currently supported file types can be found here https://github.com/jwtly10/varip/blob/main/constants.go (.env*, *.json, *.properties, *.yml, *.yaml)

More file types can easily be added by simply defining a parser, and adding the filetype to the allowed list.

Some common dependency directorys are hidden and will not be parsed for config files.

## Installation

varip can be installed on your system using one of the following methods:

### Pre-built Binaries
Pre-built binaries for varip are available under [GitHub Releases](https://github.com/jwtly10/varip/releases/). You can download the appropriate version for your operating system and architecture.

1. Navigate to the Releases page and download the latest version for your OS.
2. Unzip the downloaded file to extract the varip binary.
3. Move the varip binary to a directory in your PATH to make it globally accessible. For example on Mac/Linux:

```sh
mv varip /usr/local/bin/varip
```

<!-- ### Homebrew (macOS and Linux)
If you are on macOS or Linux, you can install varip using Homebrew, a package manager that simplifies the installation and management of software.

To install varip using Homebrew, run the following command:

```sh
brew tap jwtly10/varip
brew install varip
```
This will add the custom tap for varip and install the latest version.

Verifying the Installation
After installation, you can verify that varip is correctly installed by running:

```sh
varip --help
```
This command should display the usage information for varip, indicating that the installation was successful. -->



## Usage

varip  is easy to use and can be run from the command line by specifying a path and a pattern to search for. If no path is provided, it searches in the current directory by default.

``` sh
NAME:
   varip - Searches for environment variables in files. Searches in the current directory by default.

USAGE:
   varip [options] [path] [pattern]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose      Enable verbose debug logging (default: false)
   --errors       Display errors in output, by default errors are hidden, so only matches are shown (default: false)
   --no-color     Disable colorized output, useful if performance is slow or colors not supported by your terminal (default: false)
   --show-hidden  Show hidden files and directories (default: false)
   --help, -h     show help
```

### Examples

Search for the term "DB_PASSWORD" in all supported file types in the current directory:

```sh
varip DB_PASSWORD
```

Search for 'API_KEY' in files located in a specific directory:
``` sh
varip /path/to/configs API_KEY
```

## Development

To contribute to varip, you should have a Go development environment set up. Clone the repository, make your changes, including tests if new functionality is added. Before submitting a pull request, test your changes thoroughly.

## Contributing

Bug reports and pull requests are welcome! This is meant to be a very simple utility, but open to extending for different file types or formatting.

## License

Varip is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
