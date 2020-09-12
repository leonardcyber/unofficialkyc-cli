#!/bin/bash
deploy () {
    snapcraft remote-build --build-on=amd64 && snapcraft upload --release="${1}" ./*.snap
}

if [[ ${REF} == "refs/heads/development" ]]; then
    deploy edge
elif [[ "${REF}" == "refs/heads/staging" ]]; then
    deploy beta
elif [[ "${REF}" == "refs/heads/master" ]]; then
    deploy candidate
fi
