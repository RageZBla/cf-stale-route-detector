# CF stale route detector CLI

CLI tool to detect Go router stale routes in a Cloud Foundry deployment.

See [documentation](https://docs.cloudfoundry.org/adminguide/troubleshooting-router-error-responses.html#stale-routes-fix) about the description of the issue and fix available.

## Usage 

```
Usage:
  cf-stale-route-detector [OPTIONS] detect [detect-OPTIONS]

This command detects stale route using gorouter routing table and Diego actual LRPs exports.

Application Options:
  -v, --version            prints the om release version

Help Options:
  -h, --help               Show this help message

[detect command options]
          --routing-table= Gorouter routing table export
          --actual-lrps=   Diego actual LRPS export
          --desired-lrps=  Diego desired LRPS export
          --verbose        print details about stale route(s)
```

The routing table and actual lrps flags are required. Desired lrps is optional, if it is provided the `--verbose` flag would output additional information about the application running in diego.

The detect command exit code is `0` when no stale routes have been detected. In the case stale route have been detected the exit code is `1`.

## Gathering the required files

### Go router table

```
$ bosh ssh router -c 'sudo /var/vcap/jobs/gorouter/bin/retrieve-local-routes > /tmp/routes.json'
$ bosh scp router:/tmp/routes.json "((instance_id))"
```

### Actual and desired LRPs

```
$ bosh ssh diego_cell/0 -c '. /etc/profile.d/cfdot.sh && /var/vcap/packages/cfdot/bin/cfdot desired-lrps > /tmp/desired-lrps.json && /var/vcap/packages/cfdot/bin/cfdot actual-lrps > /tmp/actual-lrps.json'
$ bosh scp diego_cell/0:/tmp/desired-lrps.json .
$ bosh scp diego_cell/0:/tmp/actual-lrps.json .
```