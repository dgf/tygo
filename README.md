# tygo

A minimalist command-line typing test written in Go.
Measures your typing speed (WPM) and accuracy using randomized word sets
loaded from a dictionary or JSON file.

Perfect for terminal lovers, developers, or anyone who wants to practice
typing without leaving the shell.

## Features

- Fast and lightweight (compiled Go binary)
- Load custom word lists from a JSON file
- Measures **Words Per Minute (WPM)** and **accuracy**
- Real-time feedback with colored output

## Run it from source

```shell
go run main.go -dict german -punct -nums -count 20 -top 1000
```
