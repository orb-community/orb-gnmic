# Orb GNMIC

Orb GNMIC provison [gnmic]([https://github.com/openconfig/gnmic](https://github.com/openconfig/gnmic)) instances through a simple REST API using a policy mechanism. Each policy spin up a new `gnmic` process running the configuration provided by the policy.

## Project premises
**1. Single binary**: `orb-gnmic` embeds `gnmic` in its binary. Therefore, only one static binary is provided.

**2. No persistence**: `orb-gnmic` stores data in memory and in temporary files only. This adds a new paradigm to `gnmic` that is expected to run over a persisted config file as default.

**3. Compatibility**: `orb-gnmic` is basically a wrapper over the official `gnmic` which has not released a version `1.0` yet, i.e., breaking changes are expected. Any changes that occurs on its CLI will be reflected in this project.

## Docker Image
You can download and run using docker image:
```
docker run --net=host ghcr.io/orb-community/orb-gnmic run
```
## Command Line Interface (CLI)
Orb GNMIC allows some start up configuration that is listed below.
```sh
docker run --net=host ghcr.io/orb-community/orb-gnmic run --help

Run orb-gnmic

Usage:
  orb-gnmic run [flags]

Flags:
  -d, --debug                Enable verbose (debug level) output
  -h, --help                 help for run
  -a, --server_host string   Define REST Host (default "localhost")
  -p, --server_port uint     Define REST Port (default 10337)
```


## REST API
The default `orb-gnmic` address is `localhost:10222`. to change that you can specify host and port when starting `orb-gnimc`:
```sh
docker run --net=host ghcr.io/orb-community/orb-gnmic run -a {host} -p {port}
```

### Routes (v1)
`orb-gnmic` is aimed to be simple and straightforward. 

#### Get runtime and capabilities information

<details>
 <summary><code>GET</code> <code><b>/api/v1/status</b></code> <code>(gets orb-gnmic runtime data)</code></summary>

##### Parameters

> None

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json; charset=utf-8` | JSON data                                                           |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:10337/api/v1/status
> ```

</details>

<details>
 <summary><code>GET</code> <code><b>/api/v1/capabilities</b></code> <code>(gets otelcol-contrib capabilities)</code></summary>

##### Parameters

> None

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json; charset=utf-8` | JSON data                                                           |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:10337/api/v1/capabilities
> ```

</details>

#### Policies Management

<details>
 <summary><code>GET</code> <code><b>/api/v1/policies</b></code> <code>(gets all existing policies)</code></summary>

##### Parameters

> None

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json; charset=utf-8` | JSON array containing all applied policy names                      |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:10337/api/v1/policies
> ```

</details>


<details>
 <summary><code>POST</code> <code><b>/api/v1/policies</b></code> <code>(Creates a new policy)</code></summary>

##### Parameters

> | name      |  type     | data type               | description                                                           |
> |-----------|-----------|-------------------------|-----------------------------------------------------------------------|
> | None      |  required | YAML object             | yaml format specified in [Policy RFC](#policy-rfc-v1)                 |
 

##### Responses

> | http code     | content-type                       | response                                                            |
> |---------------|------------------------------------|---------------------------------------------------------------------|
> | `201`         | `application/x-yaml; charset=UTF-8`| YAML object                                                         |
> | `400`         | `application/json; charset=UTF-8`  | `{ "message": "invalid Content-Type. Only 'application/x-yaml' is supported" }`|
> | `400`         | `application/json; charset=UTF-8`  | Any policy error                                                    |
> | `400`         | `application/json; charset=UTF-8`  | `{ "message": "only single policy allowed per request" }`           |
> | `403`         | `application/json; charset=UTF-8`  | `{ "message": "config field is required" }`                         |
> | `409`         | `application/json; charset=UTF-8`  | `{ "message": "policy already exists" }`                            |
 

##### Example cURL

> ```javascript
>  curl -X POST -H "Content-Type: application/x-yaml" --data @post.yaml http://localhost:10337/api/v1/policies
> ```

</details>

<details>
 <summary><code>GET</code> <code><b>/api/v1/policies/{policy_name}</b></code> <code>(gets information of a specific policy)</code></summary>

##### Parameters

> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> |   `policy_name`   |  required | string         | The unique policy name              |

##### Responses

> | http code     | content-type                        | response                                                            |
> |---------------|-------------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/x-yaml; charset=UTF-8` | YAML object                                                         |
> | `404`         | `application/json; charset=UTF-8`   | `{ "message": "policy not found" }`                                 |

##### Example cURL

> ```javascript
>  curl -X GET http://localhost:10337/api/v1/policies/my_policy
> ```

</details>

<details>
 <summary><code>DELETE</code> <code><b>/api/v1/policies/{policy_name}</b></code> <code>(delete a existing policy)</code></summary>

##### Parameters

> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> |   `policy_name`   |  required | string         | The unique policy name              |

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json; charset=UTF-8` | `{ "message": "my_policy was deleted" }`                            |
> | `404`         | `application/json; charset=UTF-8` | `{ "message": "policy not found" }`                                 |

##### Example cURL

> ```javascript
>  curl -X DELETE http://localhost:10337/api/v1/policies/my_policy
> ```

</details>

## Policy RFC (v1)

```yaml
my_policy:
  username: admin
  password: admin
  port: 57400
  timeout: 10s
  skip-verify: true
  encoding: json_ietf

  targets:
    leaf1:57400:
    leaf2:57400:
    leaf3:57400:
    spine1:57400:
    spine2:57400:

  subscriptions:
    srl_if_oper_state:
      paths:
        - /interface[name=ethernet-1/*]/oper-state
      mode: stream
      stream-mode: sample
      sample-interval: 5s

    srl_net_instance:
      paths:
        - /network-instance[name=*]/oper-state
      mode: stream
      stream-mode: sample
      sample-interval: 5s

    srl_if_stats:
      paths:
        - /interface[name=ethernet-1/*]/statistics
      mode: stream
      stream-mode: sample
      sample-interval: 5s

    srl_if_traffic_rate:
      paths:
        - /interface[name=ethernet-1/*]/traffic-rate
      mode: stream
      stream-mode: sample
      sample-interval: 5s

    srl_cpu:
      paths:
        - /platform/control[slot=*]/cpu[index=all]/total
      mode: stream
      stream-mode: sample
      sample-interval: 5s

    srl_mem:
      paths:
        - /platform/control[slot=*]/memory
      mode: stream
      stream-mode: sample
      sample-interval: 5s

    srl_bgp_stats:
      paths:
        - /network-instance[name=*]/protocols/bgp/statistics
      mode: stream
      stream-mode: sample
      sample-interval: 5s

    srl_ipv4_routes:
      paths:
        - /network-instance[name=*]/route-table/ipv4-unicast/statistics/
      mode: stream
      stream-mode: sample
      sample-interval: 5s

    srl_ipv6_routes:
      paths:
        - /network-instance[name=*]/route-table/ipv6-unicast/statistics/
      mode: stream
      stream-mode: sample
      sample-interval: 5s

    srl_apps:
      paths:
        - /system/app-management/application[name=*]
      mode: stream
      stream-mode: sample
      sample-interval: 5s

  outputs:
    prom:
      type: prometheus
      listen: :9273
      path: /metrics
      metric-prefix: gnmic
      append-subscription-name: true
      export-timestamps: true
      debug: false
      event-processors:
        - trim-prefixes
        - up-down-map

  processors:
    trim-prefixes:
      event-strings:
        value-names:
          - ".*"
        transforms:
          - path-base:
              apply-on: "name"
    up-down-map:
      event-strings:
        value-names:
          - oper-state
        transforms:
          - replace:
              apply-on: "value"
              old: "up"
              new: "1"
          - replace:
              apply-on: "value"
              old: "down"
              new: "0"
```
