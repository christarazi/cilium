#!/bin/bash

# Source: https://kind.sigs.k8s.io/docs/user/local-registry/

source $(dirname $(readlink -ne $BASH_SOURCE))/../release/lib/common.sh
source $(dirname $(readlink -ne $BASH_SOURCE))/../backporting/common.sh

PROG=${0}
CLUSTER_NAME="${1:-}"
IMAGE="${2:-}"  # Controls the K8s version as well (e.g. kindest/node:v1.11.10)

common::argc_validate 2

# Cleanup from previous kind cluster creations
docker ps --filter=network=kind --quiet | \
    xargs --no-run-if-empty docker container rm --force
docker network rm kind || true
docker container rm --force kind-registry || true

reg_name='kind-registry'
reg_port='5000'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [[ "${running}" != "true" ]]; then
  docker run \
    -d --restart=always -p "${reg_port}:5000" --name "${reg_name}" \
    registry:2
fi

kind_cmd="kind create cluster"

if [[ -n "${CLUSTER_NAME}" ]]; then
  kind_cmd+=" --name ${CLUSTER_NAME}"
fi
if [[ -n "${IMAGE}" ]]; then
  kind_cmd+=" --image ${IMAGE}"
fi

# create a cluster with the local registry enabled in containerd
cat <<EOF | ${kind_cmd} --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
networking:
  disableDefaultCNI: true
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
    endpoint = ["http://${reg_name}:${reg_port}"]
EOF

docker network connect "kind" "${reg_name}"

for node in $(kind get nodes); do
  kubectl annotate node "${node}" "kind.x-k8s.io/registry=localhost:${reg_port}";
done

kubectl taint nodes --all node-role.kubernetes.io/master-
