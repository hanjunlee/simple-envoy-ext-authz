# Simple envoy ext_authz

![Build](https://github.com/hanjunlee/simple-envoy-ext-authz/workflows/Build/badge.svg)

## Usage

Only the `token` which passed by the argument is verified by the authorization server.

```
Usage of server
  -grpcport string
    	grpcport. (default ":50051")
  -token string
    	authenticated token. (default "foo")
```

