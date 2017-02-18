# commuter

Get commute times on the command line!

- Get commuting time between locations.
- Name and add your frequent locations for easier access.

## Install

If you have a working Go installation, you can install `commuter` like so:

```sh
$ go get github.com/KyleBanks/commuter
```

Otherwise, download `commuter` from the [Releases](https://github.com/KyleBanks/commuter/releases) page.

Run `commuter -help` to ensure you have properly installed the application.

## Usage

The first time you run `commuter`, you'll be prompted to provide a [Google Maps API Key](https://developers.google.com/console) and default location. 

**Important:** Ensure you enable the `Google Maps Distance Matrix API` at [https://developers.google.com/console](https://developers.google.com/console)!

```sh
$ commuter
(https://developers.google.com/maps/documentation/geocoding/get-api-key)
> Enter Google Maps API Key: 123APIKEY456
> Enter Default Location: 123 Main St. Toronto, Ontario
```

The API key and default location will be stored locally, and are never sent to any remote services aside from the official Google Maps API. The default location is then used by default when a `--from` location is not provided.

Next, request your commute time:

```sh
$ commuter --to "321 Maple Ave. Toronto, Ontario"
> 32 Minutes
```

If you want a commute time beginning somewhere other than your default location, you can use the `--from` flag:

```sh
$ commuter --from "123 Main St. Toronto, Ontario" --to "321 Maple Ave. Toronto, Ontario"
> 32 Minutes
```

You can also add names for your frequent locations like so:

```sh
$ commuter add --name home --address "123 Main St. Toronto, Ontario"
$ commuter add --name work --address "321 Maple Ave. Toronto, Ontario"
```

And use them as the `from` and/or `to` address:

```sh
$ commuter --from home --to work
> 32 Minutes
```

## Roadmap

- Use geolocation to allow requesting commute time from current location (ex. `commuter -c --to 123 Main St.`) where `-c` denotes to use the current location.
- Allow for alternative modes of transportation (default is driving) via `--drive`, `--walk`, `--bike`, and `--transit` flags.

