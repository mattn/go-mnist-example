# go-mnist-example

Build Docker image using [go-tflite](https://github.com/mattn/go-tflite).

## Usage

Build docker image
```
$ docker build -t mnist .
```

Test the image

```
$ docker run -i --rm mnist < 1.png
stdin is 1

$ docker run --rm -v .:/data/ mnist -f /data/4.png
/data/4.png is 4
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
