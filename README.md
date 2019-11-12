# New Relic NXLog output plugin

## Disclaimer
*IMPORTANT*: This plugin is not officially supported by New Relic and
is offered as-is.

## Limitations
*IMPORTANT*: This plugin does nothing to parse logs. Logs must be sent
in a JSON format conforming to the below structure
```
{
    timestamp: "A TIMESTAMP",
    message: "A STRING",
    attributes: {
        "attr1": "A STRING",
        "attrN": "A STRING"
    }
}
```
This can be achieved using a combination of xm_json, xm_rewrite modules
provided by NXLog.

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