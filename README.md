## Dynamic DNS Updater

A tiny go app that helps update dynamic DNS records. Currently supported DNS provider(s):

1. Namecheap.

## What does it do

In its simplest form all those lines of Go code basically just do the same as the following bash one-liner:

```console
DDNS_ip="$(curl ifconfig.me/ip 2>/dev/null)" && \
  curl "https://dynamicdns.park-your-domain.com/update?host=${DDNS_host}&domain=${DDNS_domain}&password=${DDNS_pwd}&ip=${DDNS_ip}"
```

so yeah, Go programs are way too verbose, but we digress... Some time in the future this code could start supporting 
other providers, at which point the bash one-liner won't be able to scale and the Go code will be equivalent to...
most likely: a bash script with still significantly less lines of code :)

But it's fun to write Go apps that can easily be cross-compiled for different architectures and run anywhere as a single
binary, so...

## Usage.

Set the following environmental variables: `DDNS_domain`, `DDNS_host`, `DDNS_pwd` and execute the go app in that environment.

## License 

MIT License

