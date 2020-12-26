iglu is a simple and portable program for generating Monero
cold wallets. It is designed to be used on a live operating
system, on an offline computer. The generated secret keys
should immediately be either encrypted or split (e.g. with
[gfsplit](https://git.gitano.org.uk/libgfshare.git/)) and stored on a
backup medium (I'd recommend CDs/DVDs).

I created this tool because I found the official `monero-wallet-cli` to
be cumbersome, when I just want to generate a cold wallet. I also wanted
to create a cold wallet generator, which has so little code, that most
programmers could understand it in an afternoon and thus don't have to
trust the author.

Disclaimer: I am no cryptographer, no one has audited this code and I
have not audited the code of the libraries used. This is why I recommend
to only use this tool on an offline computer with a live operating
system - this way no information can leak.

# Installation
```shell
git clone 'https://github.com/codesoap/iglu.git'
cd iglu
go build
# The iglu binary is now available at ./iglu.
# You could also install to ~/go/bin/ by executing "go install".
```

# Usage
```console
$ # Generate a cold wallet:
$ iglu
secret spend key: 29781966850873b9ec183034adf768ec3bcec4fb4cc1f4fa1297e8225fedb107
secret view key : 6bab3f01864689e642e50a3545a8acce811326c989dc0abf3c9c5b6099c4540d
primary address : 43msgHu241y1nUxfpnhWiq89V93BQd6CoJgmShqk5j2L7csCM5oajPXP9KiN1CZHgbH1BTewKBwLLBv5Fd1ZLRj13ZcvHGH

$ # Generate some subaddresses as well:
$ iglu -s 5
secret spend key: 5fd25f12c7fed7aee8d28482f1416a1294ee4fe6fd2bb21188e98d99103cfb00
secret view key : 3b7c7c44dcbd325fbaa853f981f902f7e5796f910ffa33b09d5bef1d1926ef05
primary address : 42CG9g8gJDCSUpGJMTs3ysDiZXhZaM7bN55acbSXryCM5zPDdrj5pypGnU1y228k477PdTG2AJ2gtT6zp1X3wcmf8MWQzDo
subaddress #0001: 87rnnAk44w7GSyMi3SmhpxL7t27gULE5gThkE1sxBgZ3fnRbKiTqmzL1nNQCdX8e3MRjrcE6pkkLkCeZ6USmVvwrSx19HGD
subaddress #0002: 833eUGbvtiNdGU65KVF8fzLVuCih9CUgk3pUnwTaRNW29j534hMLXTc63KguHPrJpXBfE94SRJ2PoB19QNbQsF5i5Tm99Qk
subaddress #0003: 88GyHdbWhSbBNVzLeCXzWWa6C5zvVP9uLHf5FoPQYCR7Cuug1UKHv5ibxBteEcXaQeawRLjScjxMKPopWvERWrEo76MZgLP
subaddress #0004: 89Q5CWE3dQT3CoDN9CX5f5aZYojBgzHzsRHuHSd3GvxRbtQ45Soy7mhBBbZeMs28An2QavKYHD1uDdWfKsyiAXNxVERFvvJ
subaddress #0005: 834e6PYwzzvbxuHLU5GeBZ3biGVbT4d1gfBdomX8aqFw4AujV2WjZC3ND8EHhNLuMZSNzZxPwYSQgcoiKq1KfRDTMewxnmg
```

Once you want to make your cold wallet a hot one, import the secret
spend key like this:
`./monero-wallet-cli --generate-from-spend-key <wallet-name>`.

# Where is the mnemonic?
In order to keep the code simple and because I don't find it necessary
for a cold wallet, I have decided not to generate mnemonics here. Maybe
I'll implement a separate tool for converting the secret spend key into
a mnemonic some day.
