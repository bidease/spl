# Servers.com portal

Support features of https://portal.servers.com:

- servers list
- detail information about the server

## Conf

Default path to config is `~/.spl.yml`.

```yaml
email: <email>
token: <token>
```

## Docs

```
$ curl 'https://portal.servers.com/rest/hosts' -H 'Content-Type: application/json' -H 'X-User-Email: EMAIL' -H 'X-User-Token: TOKEN' | jq
...

$ curl 'https://portal.servers.com/rest/hosts/26554' -H 'Content-Type: application/json' -H 'X-User-Email: EMAIL' -H 'X-User-Token: TOKEN' | jq
...
```
