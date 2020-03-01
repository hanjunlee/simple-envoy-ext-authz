# Simple envoy ext_authz

![Build](https://github.com/hanjunlee/simple-envoy-ext-authz/workflows/Build/badge.svg)

## What's covered?

Envoy proxy support the external authorization i.e `envoy.ext_authz` in http filters. It'll cover you can mock the external authorization server when you have a test with Envoy.

## Features

* Authorized token
* Additional Header in the response

## Usage

### Docker

```shell
$ docker run -p 50051:50051  ext-authz:0.1.0 -H foo=bar -H far=baz
# run server
2020/03/01 15:19:38 Starting gRPC Server at :50051, valid token: foo, headers: map[far:baz foo:bar].
```

### Envoy

See [envoy](./examples/envoy) in the examples directory.
