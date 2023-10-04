# Escape repeater

A simple repeater that can be used to allow escape to connect to some APIs.

## Requirements

You need to have an [Escape](https://escape.tech) account.

Before installing the repeater, you need to retrieve your repeater id.

- `ESCAPE_REPEATER_ID` : Your repeater id, read the [documentation](https://docs.escape.tech/enterprise/repeater) to know how to retrieve it.

## Install

You have multiple options to install the repeater:

- Docker image (covered in this tutorial)
- From source ([go to releases](https://github.com/Escape-Technologies/repeater/releases/latest))

You now need to run the repeater with the following environment variables:

- `ESCAPE_REPEATER_ID`: Your repeater id.

```bash
docker run -it --rm --name escape-repeater \
    -e ESCAPE_REPEATER_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx \
    escapetech/repeater:latest
```

You can find in the example folder more deployment examples.
Feel free to contribute and add your own.

## Usage

You can now go to the escape documentation follow the steps to [use your repeater](https://docs.escape.tech/enterprise/repeater).

## Contributing

Feel free to contribute to this project by creating a pull request.

If you want to test it locally, you can use the `ESCAPE_REPEATER_URL` to point to another repeater.
