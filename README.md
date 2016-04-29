# Scytale

__NOTE: This is just a toy project. You shouldn't try to actually use this
yet.__

Simple image hosting for the paranoid.

So how does this work? Well, before anything happens, the file you selected is
encrypted locally with a random key. After that, the encrypted file is uploaded
to the server and the server will respond with a url containing the location of
the file. Locally, the key is appended to that url as a location hash property.
This makes sure the key is never sent to the server if someone clicks on the
link, while the client-side JavaScript code can still read it (which it has to,
in order to decrypt the file if someone wants to view the image you linked).

Enjoy your placebo software.

## Specs

- Client-side encryption
([JavaScript crypto, fuck yeah](https://github.com/jedisct1/libsodium.js))
- XSalsa20 + Poly1305

## Requirements

- Go 1.5 or newer (for vendoring support)

## License

The entire codebase is licensed under [AGPL](LICENSE) unless stated otherwise.

## FAQ
#### Do I have to trust the server?

###### If you use the web interface

Yes, you do. The server could suddenly start serving a broken version of the
NaCl library or malicious JavaScript code that steals the key for all you know.

###### If you only use the CLI client

To not delete your files? Yes. Other than that? No.

#### Why are you using Salsa20, a stream cipher, to encrypt files?

Because for small files, it really doesn't matter whether you're using a block
cipher or a stream cipher. Salsa20 has a really good reputation and is a lot
faster than some other ciphers like AES.

Tests done with a 500 KB image
(using [SJCL](https://bitwiseshiftleft.github.io/sjcl/) for the AES
implementation):

```
salsa20 x 39.49 ops/sec ±0.89% (40 runs sampled)
aes-256 x 3.66 ops/sec ±3.44% (14 runs sampled)
```

#### What about all those horror stories I heard about JavaScript crypto?

They're all true. You can use the included desktop [CLI client](client) if you
prefer.

#### The html/css/whatever looks horrible. What's up with that?

Tell me about it.
