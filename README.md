# commuter

[![Build Status](https://travis-ci.org/KyleBanks/commuter.svg?branch=master)](https://travis-ci.org/KyleBanks/commuter) &nbsp;
[![GoDoc](https://godoc.org/github.com/KyleBanks/commuter?status.svg)](https://godoc.org/github.com/KyleBanks/commuter) &nbsp;
[![Go Report Card](https://goreportcard.com/badge/github.com/KyleBanks/commuter)](https://goreportcard.com/report/github.com/KyleBanks/commuter)

Get commute times on the command line!

- Get commuting time between locations.
- Name and add your frequent locations for easier access.

## Install

Download the appropriate `commuter` binary from the [Releases](https://github.com/KyleBanks/commuter/releases) page.

Alternatively, if you have a working Go installation, you can install `commuter` like so:

```sh
$ go get github.com/KyleBanks/commuter
```

## Usage

The first time you run `commuter`, you'll be prompted to provide a [Google Maps API Key](https://developers.google.com/console) and default location. 

**Important:** Ensure you enable the `Google Maps Distance Matrix API` at [https://developers.google.com/console](https://developers.google.com/console)!

```sh
$ commuter
> Enter Google Maps API Key: (developers.google.com/console)
123APIKEY456
> Enter Your Default Location: (ex. 123 Main St. Toronto, Canada)
123 Main St. Toronto, Ontario
```

The API key and default location will be stored locally, and are never sent to any remote services aside from the official Google Maps API. The default location is then used by default when a `-from` or `-to` location is not provided.

Next, request your commute time:

```sh
# From your default to a specific location:
$ commuter -to "321 Maple Ave. Toronto, Ontario"
> 32 Minutes

# From a specific location to your default:
$ commuter -from "Toronto, Ontario"
> 20 Minutes
```

If you want a commute time beginning and ending somewhere other than your default location, you can use supply full locations for both the `-from` and `-to` flags:

```sh
$ commuter -from "123 Main St. Toronto, Ontario" -to "321 Maple Ave. Toronto, Ontario"
> 32 Minutes
```

You can also add names for your frequent locations like so:

```sh
$ commuter add -name home -location "123 Main St. Toronto, Ontario"
$ commuter add -name work -address "321 Maple Ave. Toronto, Ontario"
```

And use them as the `from` and/or `to` location:

```sh
$ commuter -from home -to work
> 32 Minutes
```

## License

```
MIT License

Copyright (c) 2017 Kyle Banks

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```