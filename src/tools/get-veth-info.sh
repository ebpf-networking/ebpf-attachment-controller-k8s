#!/bin/bash

# Based on - https://stackoverflow.com/questions/21724225/docker-how-to-get-veth-bridge-interface-pair-easily/28613516#28613516

function containerid_to_veth_containerd() {

    if [ -z "${CONTAINER_ID}" ]
    then
        printf '{"error":"container id must be passed as first argument"}'
        exit 1
    fi

    # Get the container network namespace from container id.
    local proc_namespace=$(ctr --namespace k8s.io c info ${CONTAINER_ID} | jq '.Spec.linux.namespaces[] | select(.type == "network") | .path')
    local pid=$(echo ${proc_namespace} | awk -F / '{print $3}')

    # Make the container's network namespace available to the ip-netns command:
    local id=$(echo ${CONTAINER_ID} | cut -c1-5)
    mkdir -p "/var/run/netns"
    ln -sf "/proc/${pid}/ns/net" "/var/run/netns/${id}"

    # Get the vpeer full name.
    local vpeer_name=$(ip netns exec "$id" ip --brief link show eth0 | head -n1 | awk '{print $1}')

    # Get the veth index associated with vpeer
    veth_index=$(echo ${vpeer_name} | cut -d '@' -f 2 | sed -e "s/^if//")

    # Find the veth name
    veth_name=$(ip link show | grep "^${veth_index}:" | sed "s/${veth_index}: \(.*\):.*/\1/" | cut -d '@' -f 1)

    # Find the veth mac
    veth_mac=$(ip link show ${veth_name} | sed -n 2p | awk '{print $2}')

    # Clean up the netns symlink, since we don't need it anymore
    rm -f "/var/run/netns/${id}"
}

CONTAINER_ID="${1}"
containerid_to_veth_containerd
printf '{"veth-name":"%s","veth-id":"%s","veth-mac":"%s"}\n' "$veth_name" "$veth_index" "$veth_mac"