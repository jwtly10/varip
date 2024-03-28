# varip

A lightweight and fast command line environment variables ripper (think simple Grep for config files).

Notes: 
1. Currently supported file types can be found here https://github.com/jwtly10/varip/blob/main/constants.go (.env*, *.json, *.properties, *.yml, *.yaml).
2. More file types can easily be added by simply defining a parser, and adding the filetype to the allowed list.
3. Some common dependency directorys are hidden and will not be parsed for config files.
4. Incorrectly formatted files throw errors and are skipped. Only valid files will be parsed.

## Installation

varip can be installed on your system using one of the following methods:

### Pre-built Binaries

Pre-built binaries for varip are available under [GitHub Releases](https://github.com/jwtly10/varip/releases/). You can download the appropriate version for your operating system and architecture. However note that they are NOT signed, so your respective OS may refuse to run them. 

I will not provide steps to circumvent this, but feel free to do so.

Instead, it's easier to just clone and build the app. Assuming you have git and Go 1.22+ installed, you can install varip using the following script:
### Linux/Mac

Ensure you are in a directory where you want to clone varip.

This script will install varip to /usr/local/bin, which is on the PATH so allows access to `varip` anywhere. Feel free to customize it to your needs.

``` bash
#!/bin/bash

echo "Cloning the varip repository..."
git clone https://github.com/jwtly10/varip &&

echo "Changing directory to varip..."
cd ./varip &&

echo "Building the varip binary..."
go build &&

echo "Moving the varip binary to /usr/local/bin..."
sudo mv varip /usr/local/bin/varip

echo "Installation complete. You can now use varip by typing 'varip' in your terminal."

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

There are some abstractions for the sake of a clean interface.
1. Errors during file parsing are hidden by default. To display errors, use the `--errors` flag.
2. Verbose debug logging is disabled by default. To enable verbose debug logging, use the `--verbose` flag.


This leads to a cleaner interface of just seeing the matches you are looking for.


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
