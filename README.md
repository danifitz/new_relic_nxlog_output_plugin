# New Relic NXLog output plugin

## Usage
Follow instructions in `Build` to create a binary.

Use the binary with the NXLog output module `om_exec`

https://nxlog.co/documentation/nxlog-user-guide/om_exec.html

```
<Output someprog>
    Module  om_exec
    Command /path/to/binary
    Arg     --NEW_RELIC_LICENSE_KEY <YOUR_KEY>
    Arg     --NEW_RELIC_LOG_URI DEFAULTS_TO_https://log-api.newrelic.com/log/v1
</Output>
```

## Build
Run `go build` to generate a binary