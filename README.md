### Info

Encrypt files fast and secure. Support linux and windows.

Based on [cry](https://github.com/wille/cry), with more function.

### usage

mode genlocal: first use，gen priv key to encrypt file and store key locally.

mode local: use local key to encrypt or decrypt, don't update key.

mode remote: use remote key to encrypt or decrypt, don't update key.

mode genremote: first use，gen priv key to encrypt file and store key remotely.

>If files are already encrypted, it won't encrypt anymore, but with gen modes the key still get stored or uploaded.

### todo

1. encrypt big file. 内存占用问题。