![Commuter](./_misc/commuter.png)

[![Build Status](https://travis-ci.org/KyleBanks/commuter.svg?branch=master)](https://travis-ci.org/KyleBanks/commuter) &nbsp;
[![GoDoc](https://godoc.org/github.com/KyleBanks/commuter?status.svg)](https://godoc.org/github.com/KyleBanks/commuter) &nbsp;
[![Go Report Card](https://goreportcard.com/badge/github.com/KyleBanks/commuter)](https://goreportcard.com/report/github.com/KyleBanks/commuter)

Get commute times on the command line!

- Get commuting time between locations.
- Name and add your frequent locations for easier access.
- Geolocation allows you to use your current location without typing an address.

## Install

Download the appropriate `commuter` binary from the [Releases](https://github.com/KyleBanks/commuter/releases) page.

Alternatively, if you have a working Go installation, you can install `commuter` like so:

```sh
$ go get github.com/KyleBanks/commuter
```

## Usage

The first time you run `commuter`, you'll be prompted to provide a [Google Maps API Key](https://developers.google.com/console) and default location. 

**Important:** Ensure you enable the *Google Maps Distance Matrix API* at [developers.google.com/console](https://developers.google.com/console). If you want to use the `-from-current` and `-to-current` flags, you will also need to enable the *Google Maps Geolocation API*.

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
32 Minutes

# From a specific location to your default:
$ commuter -from "Toronto, Ontario"
20 Minutes
```

If you want a commute time beginning and ending somewhere other than your default location, you can use supply full locations for both the `-from` and `-to` flags:

```sh
$ commuter -from "123 Main St. Toronto, Ontario" -to "321 Maple Ave. Toronto, Ontario"
32 Minutes
```

### `commuter add`

You can also add names for your frequent locations like so:

```sh
$ commuter add -name home -location "123 Main St. Toronto, Ontario"
$ commuter add -name work -address "321 Maple Ave. Toronto, Ontario"
```

And use them as the `from` and/or `to` location:

```sh
$ commuter -from home -to work
32 Minutes
```

### `commuter list`

To see a list of all your named locations:

```sh
$ commuter list
default: 123 Main St. Toronto, Ontario
    gym: 1024 Fitness Lane Toronto, Ontario
   work: 321 Maple Ave. Toronto, Ontario
```

### Using Your Current Location

If you [enabled](https://developers.google.com/console) the *Google Maps Geolocation API* for your API key, you can use the `-from-current` and `-to-current` flags to use your current location. This is done by attempting to use your IP Address to determine your latitude and longitude, and use that as either the start or destination of your commute:

```sh
$ commuter -from-current -to work
32 Minutes
$ commuter -from gym -to-current
12 Minutes
```

### Travel Modes

By default, `commuter` assumes you are driving between locations. However, you can specify one or more commute methods using the `-drive`, `-walk`, `-bike` and `-transit` flags, like so:

```sh
# Single mode:
$ commuter -walk -to work
7 Hours 50 Minutes

# Multiple modes:
$ commuter -walk -transit -drive -bike -to work
Drive: 30 Minutes
Walk: 7 Hours 50 Minutes
Bike: 2 Hours 45 Minutes
Transit: 1 Hour 17 Minutes
```

And of course the different travel modes can be combined with your current location:

```sh
$ commuter -bike -from-current -to gym
2 Hours 18 Minutes
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
