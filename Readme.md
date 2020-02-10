
# impersonate(1)

[Auth0](https://auth0.com/) user impersonation utility for when you need to see
exactly what your customers see.

## Installation

From [gobinaries.com](https://gobinaries.com):

```sh
$ curl -sf https://gobinaries.com/tj/impersonate | sh
```

From source:

```sh
$ go get github.com/tj/impersonate
```

## Setup

You'll need to set the following "global" (account level) env vars:

```sh
export AUTH0_CLIENT_ID=xxxxxxx
export AUTH0_CLIENT_SECRET=xxxxxxx
```

## Usage

To impersonate a user pass the Client ID of your application (not your account),
the "impersonator" user ID (you), and the ID of the user.

```sh
$ impersonate --account apex-inc --client-id xxxxx --impersonator-id 'github|yyy' 'github|zzz'
```

If this gets annoying or you have multiple applications, you may want to alias in your profile:

```sh
alias impersonate_myapp="impersonate --account apex-inc --client-id xxxxx --impersonator-id 'github|yyy'"
```

Then all you need is:

```sh
$ impersonate_myapp 'github|zzz'
```

## Badges

![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)
[![](http://apex.sh/images/badge.svg)](https://apex.sh/ping/)

---

> [tjholowaychuk.com](http://tjholowaychuk.com) &nbsp;&middot;&nbsp;
> GitHub [@tj](https://github.com/tj) &nbsp;&middot;&nbsp;
> Twitter [@tjholowaychuk](https://twitter.com/tjholowaychuk)
