#!/bin/bash

# Note:-This script assumes vpeer_name as eth0 inside container. 

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
    netns_id=$(echo ${CONTAINER_ID} | cut -c1-5)

    mkdir -p "/var/run/netns"
    ln -sf "/proc/${pid}/ns/net" "/var/run/netns/${netns_id}"

    vpeer_name="eth0" #$(echo ${vpeer_fullname} | cut -d '@' -f 1)

    # Get the vpeer full name.
    local vpeer_fullname=$(ip netns exec "$netns_id" ip --brief link show ${vpeer_name} | head -n1 | awk '{print $1}')

    # Get vpeer index
    vpeer_index=$(ip netns exec "$netns_id" ip link show ${vpeer_name} | head -n1 |  awk '{print $1}' | sed -e "s/\://")

    # Find the vpeer mac
    vpeer_mac=$(ip netns exec "$netns_id" ip link show ${vpeer_name} | sed -n 2p | awk '{print $2}')

    # Get the veth index associated with vpeer
    veth_index=$(echo ${vpeer_fullname} | cut -d '@' -f 2 | sed -e "s/^if//")

    # Find the veth name
    veth_name=$(ip link show | grep "^${veth_index}:" | sed "s/${veth_index}: \(.*\):.*/\1/" | cut -d '@' -f 1)

    # Find the veth mac
    veth_mac=$(ip link show ${veth_name} | sed -n 2p | awk '{print $2}')

    # Clean up the netns symlink, since we don't need it anymore
    # rm -f "/var/run/netns/${id}"
}

CONTAINER_ID="${1}"
containerid_to_veth_containerd
printf '{"veth-name":"%s","veth-id":"%s","veth-mac":"%s","vpeer-name":"%s","vpeer-id":"%s","vpeer-mac":"%s","netns":"%s"}\n' \
 "$veth_name" "$veth_index" "$veth_mac" "$vpeer_name" "$vpeer_index" "$vpeer_mac" "$netns_id"