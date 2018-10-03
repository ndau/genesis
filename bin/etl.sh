#!/bin/bash

dependencies=(go glide grep jq toml2json json2toml)
for tool in "${dependencies[@]}"; do
    if ! command -v "$tool" > /dev/null  ; then
        >&2 echo  "This script depends on $tool. Install it and try again."
        exit 1
    fi
done

root_path=$(git rev-parse --show-toplevel)
cd "$root_path" || exit 1
if [ ! -x ./etl ]; then
    glide install
    go build ./cmd/etl
fi

if [ -z "$NDAUHOME" ]; then
    NDAUHOME="$HOME/.ndau"
fi

config="$NDAUHOME"/ndau/config.toml
if [ ! -f "$config" ]; then
    # shellcheck disable=SC2016
    >&2 echo 'Could not locate ndau configuration file. Check $NDAUHOME'
    >&2 echo "Expected it at $config"
    exit 1
fi

config_bk="$config".bk
if grep -q UseMock "$config"; then
    # we probably need to adjust the config file temporarily:
    # chances are good that it's configured within a docker context, which
    # means that the mockfile isn't quite where it's claimed to be

    mockfile="$NDAUHOME"/ndau/mock-chaos.msgp
    if [ ! -f "$mockfile" ]; then
        # shellcheck disable=SC2016
        >&2 echo 'Could not locate mockfile. Check $NDAUHOME'
        exit 1
    fi

    >&2 echo 'Adjusting ndau configuration for local mockfile...'

    cp "$config" "$config_bk"

    toml2json "$config_bk" |\
    jq ".UseMock=\"$mockfile\"" |\
    json2toml > "$config"
fi

./etl

if [ -f "$config_bk" ]; then
    >&2 echo 'Resetting configuration...'
    mv "$config_bk" "$config"
fi

>&2 echo "updating empty app hash..."
cd ../ndau || exit 1
bin/update-hash.sh
