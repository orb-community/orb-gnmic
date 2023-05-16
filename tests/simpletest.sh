#!/usr/bin/env bash

CURRENT_DIR=$pwd
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd $SCRIPT_DIR

if [ ! -f "../build/orb-gnmic" ]; then
    make -C ../ build
fi

../build/orb-gnmic run &
PID=$!

sleep 2

ret=$(curl -o /dev/null -s -w "%{http_code}\n" -X POST --location 'localhost:10337/api/v1/policies' -H 'Content-Type: application/x-yaml' --data 'policy_test:
  config:
    username: admin
    password: admin
    port: 57400
    timeout: 10s
    skip-verify: true
    encoding: json_ietf
    log: true
    targets: "leaf1:57400"
    subscriptions:
      srl_if_oper_state:
        paths:
          - "/interface[name=ethernet-1/*]/oper-state"
        mode: stream
        stream-mode: sample
        sample-interval: 5s
    outputs:
      prom:
        type: prometheus
        listen: "127.0.0.1:9273"
        path: /metrics
        metric-prefix: gnmic
        append-subscription-name: true
        export-timestamps: true
        debug: false
        event-processors:
          - trim-prefixes
          - up-down-map')

#kill $PID
cd $CURRENT_DIR

echo $ret

if [[ $ret != 201 ]]; then
  exit 1
fi

exit 0

