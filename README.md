<div align="center">
	<p>
		<a href="https://slack.linodians.com">
			<img alt="Cloud Unpacked" src="img/logo-badge-circle.svg" width="75" />
		</a>
	</p>
	<h1>Cloud Unpacked - Dynamic DNS</h1>
	<h3>a Dynamic DNS client for VPS cloud providers</h3>
</div>

[![Build Status](https://circleci.com/gh/cloud-unpacked/cu-ddns.svg?style=shield)](https://circleci.com/gh/cloud-unpacked/cu-ddns) [![Software License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/cloud-unpacked/cu-ddns/master/LICENSE) [![Follow @CloudUnpacked](https://img.shields.io/twitter/follow/CloudUnpacked.svg?label=Follow%20@CloudUnpacked)](https://twitter.com/intent/follow?screen_name=CloudUnpacked)

***This project is brand new and in alpha right now. You may notice a lack of polish and/or drastic changes.***

`cu-ddns` is a dynamic DNS client that uses VPS cloud providers such as Linode for DNS.
This tool allows pointing a DNS hostname such as `home.example.com` to an IP address that may change regularly.
The typical scenario is having a domain name point to your home IP address however those that travel a lot would find it useful as well.


## Table of Contents

- Installing
- Configuring
- Features


## Providers

`cu-ddns` supports the following providers:

- Linode DNS
- Cloudflare DNS


## Installing

### Debian Package (.deb) Instructions

Download the `.deb` file to the desired system.

For graphical systems, you can download it from the [GitHub Releases page][gh-releases].
Many distros allow you to double-click the file to install.
Via terminal, you can do the following:

```bash
wget https://github.com/cloud-unpacked/cu-ddns/releases/download/v0.1.0/cu-ddns_0.1.0_amd64.deb
sudo dpkg -i cu-ddns_0.1.0_amd64.deb
```

`0.1.0` and `amd64` may need to be replaced with your desired version and CPU architecture respectively.

### Linux Snap

If you install via Snap, some commands will need to be prefixed with `sudo`.

More instructions coming soon.


## Configuring

After installation, run the `configure` command to setup the client and the `start` command to start it running.

```bash
sudo cu-ddns configure
sudo cu-ddns start
```

### Cloudflare

When creating the Cloudflare API token, the following permissions are needed: `All zones - Zone:Edit, DNS:Edit`.


## Features

*Multiple Providers* - Linode and Cloudflare are supported.
DigitalOcean DNS next on the list.

*IPv4/6 Support* - IPv4 is supported with IPv6 coming in the near future.


## License

This repository is licensed under the MIT license.
The license can be found [here](./LICENSE).



[gh-releases]: https://github.com/cloud-unpacked/cu-ddns/releases
