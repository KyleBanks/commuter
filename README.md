# commuter

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
# Default to a location
$ commuter -to "321 Maple Ave. Toronto, Ontario"
> 32 Minutes

# A location to your default
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

