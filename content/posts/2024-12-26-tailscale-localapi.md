+++
title = 'Tailscale LocalAPI'
date = 2024-12-26T15:25:08+01:00
draft = false
tags = ['Tailscale', 'LocalAPI', 'TIL']
+++

## Context

I use Tailscale to secure network communication, and I wanted to authenticate users coming from the Tailscale network.

The goal was to reproduce the behavior in the [golink](https://github.com/tailscale/golink) application.
When users come from the Tailscale network, the application should authenticate them and give them access to the application based on roles defined in Tailscale ACL.

## Solution

Tailscale exposes an HTTP API through a Unix Socket named `LocalAPI`.

It can be used for many things. In my case, I needed to authenticate a user coming from Tailnet.
I used it to get information about a Tailscale user by calling the `whois` endpoint.

## How to use it

Example with `curl`:

```shell
curl \
    --unix-socket /run/tailscale/tailscaled.sock \
    'http://local-tailscaled.sock/localapi/v0/whois?addr=remoteAddress[:remotePort]'
```

**Update**: the `remotePort` is not mandatory but can be useful to get information when using Tailscale in [userspace mode](https://tailscale.com/kb/1112/userspace-networking).
Thank you [Brad Fitzpatrick](https://bsky.app/profile/bradfitz.com) for the clarification [in Bluesky](https://bsky.app/profile/bradfitz.com/post/3leabcqqvek23).

JSON returned (truncated):

```json
{
	"Node": {
		"ID": 1234567890,
		"StableID": "foo",
		"Name": "computer.tailnet-name.ts.net.",
		"User": 12345,
		"ComputedName": "computer",
		"ComputedNameWithHost": "computer (localhost)"
	},
	"UserProfile": {
		"ID": 12345,
		"LoginName": "john.doe@company.tld",
		"DisplayName": "DOE John",
	},
	"CapMap": {
		"company.tld/cap/app": [
			{
				"role": [
					"foo",
					"bar"
				]
			}
		]
	}
}
```

We have enough information to authenticate the user, and we can even map some roles in Tailscale ACL thanks to the [`Grants`](https://tailscale.com/kb/1337/acl-syntax#grants) section.

This API is equivalent to the `tailscale whois` command, but it can be used from any language without running the `tailscale` binary.

## Documentation

Unfortunately, LocalAPI is not documented yet.
A Go SDK exists to interact with it: [Tailscale LocalClient](https://pkg.go.dev/tailscale.com/client/tailscale#LocalClient).
When using another language we have to check API calls in the source code to reproduce them:

- [Go LocalClient source code](https://github.com/tailscale/tailscale/blob/v1.78.3/client/tailscale/localclient.go#L60)
- [Go LocalAPI handler](https://github.com/tailscale/tailscale/blob/v1.78.3/ipn/localapi/localapi.go#L76)

## Disclaimer

Use with caution, we have no official guarantee that this API will not change in the future due to the lack of documentation.
