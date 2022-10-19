# Remote Execute

Remotely execute commands on your server through an HTTP Post request.

## Usage

* Create a `config.json` file with the following parameters (example below):
```
{
    "ListenHost": "127.0.0.1",
    "ListenPort": 8000,
    "Tokens": [
        "BdpHt88ZsqiZibHeLpaXkJXwGUaQKi4uVcUTGnukiBADSs5ARU",
        "GrMM8BMbqM89TxZdsf6cM5oqTBFVmojFg2LcJ4w6fbSkpkoVWV"
    ]
}
```
* [Install go](https://go.dev/doc/install)
* Run the server
```
$ go build
$ ./main
```
* Visit `http://127.0.0.1:8000/` in your browser or in a terminal (emulating an HTTP POST Request)
```
$ curl -X POST -H "Authentication: Bearer <token_from_config.json>" -F "command=wget https://thejaswi.info/" http://127.0.0.1:8000/
```

## Optional

To run untrusted user input, it's best to sandbox the server using tools like [BubbleWrap](https://github.com/containers/bubblewrap), [Firejail](https://github.com/netblue30/firejail) etc.

Here's an example with bubblewrap (on Fedora) that allows only one writable directory (`/user_dir`):

```
$ mkdir /tmp/user_dir
$ bwrap --ro-bind /usr /usr \
        --ro-bind /lib64 /lib64 \
        --ro-bind /etc/ssl/certs /etc/ssl/certs \
        --ro-bind /etc/pki /etc/pki \
        --ro-bind /etc/alternatives /etc/alternatives \
        --proc /proc \
        --dev /dev \
        --ro-bind . /run/0 \
        --chdir /run/0 \
        --bind /tmp/user_dir /user_dir \
        --unshare-all --share-net \
        -- ./main
```

## License

Available under `MIT License`. Please check the `LICENSE` file for details.
