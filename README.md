# scytale

Simple file hosting for the paranoid.

## License

The entire codebase is licensed under [AGPL](LICENSE) unless explicitly stated otherwise.

## Encryption details

As mentioned in the intro, files can be encrypted client-side before uploading.
The algorithm used for this is __AES__ in __GCM__ mode with a random __256-bit__
key.

The key is included in the hash portion of URL. Web browsers don't forward this
part to the server, but can still use it in their JavaScript code to perform
decryption of the file.

The web client uses [SJCL](https://bitwiseshiftleft.github.io/sjcl/) for the AES
implementation. The CLI client uses the AES implementation in Go's standard
library (crypto/cipher and crypto/aes)

## Contributing

Contributions are greatly appreciated. If you plan on making some big changes,
please open an issue to discuss it first. I'm trying to keep this service as
simple as possible so the chance of big changes getting merged in without any
prior discussion is very low.

## Installation

TODO

### Requirements

- Go 1.5 or newer (for vendoring support)

## FAQ
#### Do I have to trust the server?

###### If you use the web interface

Yes, you do. The server could suddenly start serving a broken version of the
crypto library or malicious JavaScript code that steals the key without you
noticing.

###### If you only use the CLI client

To not delete your files? Yes. Other than that? No.

#### What about all those horror stories I heard about JavaScript crypto?

They're all true. You can use the included desktop [CLI client](client) if you
prefer.
