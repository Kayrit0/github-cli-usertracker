# UserTracker

A CLI tool for tracking GitHub users and their repositories.

## Installation

```bash
go build -o usertracker.exe
```

## Usage

### Get user information

```bash
usertracker user <username>
```

Example:
```bash
usertracker user Kayrit0
```

Displays:
- Login and name
- User ID
- Company and location
- Email and bio
- Number of public repositories
- Followers and following count
- Account creation date
- Profile URL

### Get user repositories

```bash
usertracker repos <username>
```

Example:
```bash
usertracker repos Kayrit0
```

Displays for each repository:
- Full name (owner/repo)
- Description
- Programming language
- Stars, forks, and open issues count
- Repository URL
- Last update date
- [FORK] and [PRIVATE] labels

## Requirements

- Go 1.25.5+
- Internet connection for GitHub API access

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
