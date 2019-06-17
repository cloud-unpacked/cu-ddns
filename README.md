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


## Installing

Instructions coming soon.

### Linux Snap

If you install via Snap, some commands will need to be prefixed with `sudo`.


## Configuring

After installation, run the `configure` command to setup the client and the `start` command to start it running.

```bash
cu-ddns configure
cu-ddns start
```


## Features

*Multiple Providers* - This is currently a lie.
Linode DNS is supported with DigitalOcean DNS next on the list.

*IPv4/6 Support* - both IPv4 and IPv6 are supporting.
Currently, you can't specifically choose though. Whichever is the default route/protocol of your local machine, that's what will be used.


## License

This repository is licensed under the MIT license.
The license can be found [here](./LICENSE).
