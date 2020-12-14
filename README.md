iglu is a simple and portable program for generating Monero
cold wallets. It is designed to be used on a live operating
system, on a offline computer. The generated secret keys
should immediately be either encrypted or split (e.g. with
[gfsplit](https://git.gitano.org.uk/libgfshare.git/)) and stored on a
backup medium (I'd recommend CDs/DVDs).

# Usage
```console
$ # Generate a cold wallet:
$ iglu
secret spend key: 29781966850873b9ec183034adf768ec3bcec4fb4cc1f4fa1297e8225fedb107
secret view key : 6bab3f01864689e642e50a3545a8acce811326c989dc0abf3c9c5b6099c4540d
primary address : 43msgHu241y1nUxfpnhWiq89V93BQd6CoJgmShqk5j2L7csCM5oajPXP9KiN1CZHgbH1BTewKBwLLBv5Fd1ZLRj13ZcvHGH
```
