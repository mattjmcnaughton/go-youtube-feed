# go-youtube-feed

Just a fun little coding project to generate a Atom Feed for a YouTube channel.

## Developing

```
nix develop
```

## Running locally

```
cp .env.sample .env  # Populate `YOUTBUE_API_KEY`
go run main.go [OPTIONS]  # Use --help for options
```

## Obtaining API key

1. From the Google Developer Console, create a new project (i.e.
   `go-youtube-feed`).
2. Enable the [YouTube Data API v3](https://console.cloud.google.com/apis/library/youtube.googleapis.com).
3. Create API credentials for the `go-youtube-feed` project. Use the `API
   Key` credential type.
   - Restrict the API Key to `YouTube Data API v3`.
   - If you want, restrict the API key by IP address. You can find your IP
     address at https://nordvpn.com/what-is-my-ip/.
   - Give the API key a sensible name.

## TODO

- Add an option to dump feeds for all channels someone subscribes to.
- Add unit and integration testing.
- Configure logging via a tool like `logrus`
- Write a react frontend which we can trigger via `server`.
- Write `rs-youtube-feed` to experiment w/ writing in Rustlang.
