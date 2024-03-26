package main

// ignoredDirectories contains a list of directory names that should be ignored during the search process.
// These are directories typically used for dependencies or build outputs that are not relevant to the search.
// A flag input is provided to allow the user to override this list and search in these directories.
var ignoredDirectories = []string{"vendor", "node_modules", "__pycache__", "build", "dist", ".git", "tmp"}

// supportedFileTypes lists the file extensions of files that will be searched for the specified patterns.
// This allows varip to focus on likely candidates for configuration files while skipping over unrelated file types.
var supportedFileTypes = []string{".env*", "*.json", "*.properties", "*.yml", "*.yaml"}
