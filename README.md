# Gomail
Gomail is a command line application that aims to simulate, not replace, using Gmail in the command line such as checking new emails, sending emails and searching for them.

I decided to build this project out of fun and curiosity to see if I could turn this from an idea to something that actually exists. This is also my first time using `Go` and learning about TUI libraries like [Bubbletea](https://github.com/charmbracelet/bubbletea/tree/main) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Running locally

### Requirements
1. Have [Go](https://go.dev/doc/install) (>=v1.23.0) installed on your machine
2. [Google account](https://www.google.com/intl/en-GB/account/about/)

### Clone the repository
Clone the repo by clicking on the Code button above and copying the link. Depending on your machine, open either the command line or powershell for Windows or the integrated terminal on macOS/linux and type 

```git clone <repo-url>```

replace `repo-url` with the link you copied earlier and press `Enter`. Naviagte to the `gomail` folder that was created when you cloned the repo and open it in a code editor or IDE of your choice.

### Run the app
Run `go run .` in your terminal and select one of options presented to you. After your selection, you will be asked to copy a link and paste it in your browser which will allow you to obtain a authorization code. Select your google account and you will be asked if you trust Gomail to access your account and view what it will gain access to. Select `Continue` if you wish to proceed.

Once you proceed, the authorization code will be in the URL of the current page after the `&code=` parameter. Copy what is after that up until but not including `&scope=`. Paste the code back into your terminal and press `Enter`.
