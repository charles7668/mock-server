# Mock Server

A mock backend for frontend testing.

## Usage

The default input file is named `config.json`. You can use `--file` or `-f` to specify your config file.

Write your route settings in JSON. Here is an example [config.json](./config.json):

- `port`: The port your server will listen on.
- `GET`: Method for handling `GET` requests.
  - `path`: Route path, this setting needs a value.
  - `body`: Response value. If not set, the response will be empty.
  - `status code`: HTTP status code for the response. If not set, the status code is 200.
  - `headers`: Set extra response headers. You can also leave it empty.

```json
{
  "port": 3000,
  "GET": [
    {
      "path": "/",
      "body": {
        "message": "Hello World"
      },
      "headers": {
        "Content-Type": "application/json"
      }
    },
    {
      "path": "/about",
      "status": 403,
      "body": {
        "message": "about"
      }
    }
  ]
}
```

This configuration will set up a server on port 3000 with the following routes:

- A `GET` request to `/` will return a JSON response `{"message": "Hello World"}` with a `Content-Type` header of `application/json`.
- A `GET` request to `/about` will return a status code of 403 and a JSON response `{"message": "about"}`.
