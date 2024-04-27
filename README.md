# HTTP Request Counter

HTTP Request Counter is a Go application that counts the total number of requests received during the previous 60 seconds.

## Features

- Counts total requests received in the last 60 seconds.
- Persists data to a file to maintain count across application restarts.
- Provides an HTTP server to serve the current request count.

## Usage

### Running the Application

To build and run the application, execute the following command:

```bash
make run
```

### Making Requests

To retrieve the current request count execute cURL:

```bash
curl http://localhost:8080/counter
```

### Testing

To run unit tests, execute the following command:

```bash
make test
```