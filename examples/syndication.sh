#!/bin/bash

usage="
Usage:
$(basename "$0") -f FROM -t TO [-g] [-h]  -- syndication script for Mohawk servers

where:
    -h  show this help text
    -f  from url e.g. https://metrics.node.com:8443
    -t  to url e.g. https://metrics.cluster.com:8443
    -g  update tags

example:
   $(basename "$0") -f https://localhost:8443 -t https://localhost:8444 -g"

while getopts hgf:t: option
do
 case "${option}"
 in
 h) echo "${usage}"; exit;;
 f) FROM_URL=${OPTARG};;
 t) TO_URL=${OPTARG};;
 g) TAGS="YES";;
 esac
done

if [ -z "${FROM_URL}" ]; then
  echo "${usage}"
  exit 1
fi

if [ -z "${TO_URL}" ]; then
  echo "${usage}"
  exit 1
fi

# Get data from source
curl -ks -H "Accept-Encoding: gzip"  -X GET  ${FROM_URL}/hawkular/metrics/metrics > syn.json.gz

# Post data to higher teir
curl -ks -H "Content-Encoding: gzip" -X POST ${TO_URL}/hawkular/metrics/gauges/raw --data-binary @syn.json.gz

if [ -z "${TAGS}" ]; then
  exit
fi

# Post tags to higher teir
curl -ks -H "Content-Encoding: gzip" -X PUT ${TO_URL}/hawkular/metrics/gauges/tags --data-binary @syn.json.gz
