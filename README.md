# Marvel
Marvel API in Golang

## Usage
Run the following from terminal:
```
go get github.com/jtejido/go-marvel
cd <golang-path>/go-marvel
go build
./go-marvel -config=<config_path>
```

### Configuration file
The API client reads a yaml config file with the ff properties:

```
key: <API_KEY>
secret: <API_PRIVATE_KEY>
listen: 0.0.0.0:8080  # domain it listens to
timeout: 10             # number of seconds for a request timeout to Marvel's api, if empty, it will wait indefinitely
token: <RANDOM_HASH> 	# session token, optional, can be empty
debug: true
```

When running the binary without **-config** flag, it will attempt to lookup from **MARVEL_API_CONFIG** environment, otherwise an error.

## Usage

See this [doc](https://github.com/jtejido/go-marvel/tree/master/api/docs) for fetching Character/s and Comic/s requests.

### Caching
This [General Info](https://developer.marvel.com/documentation/generalinfo) page pretty much details how one should use **ETag** for request optimization.

All requests goes thru a middleware for caching ETag per URI (URI => ETag) on an in-memory KV store, this ETag will be set upon initial request, and will be used and checked for any succeeding requests, an updated ETag would have to renew this value per URI. On initial request, the body fetched from Marvel's API will be stored on
another KV store instance (ETag => Body), when a 304 is received from Marvel for the succeeding ones, this stored body shall be used, thus avoiding the need
to re-fetch it from Marvel. An updated Etag would have to remove/reset the old data.
