#!/usr/bin/env bash

KUBERNETES_VERSION=$1
export DOCKER_BUILDKIT=1
rm -rf kubernetes
git clone https://github.com/kubernetes/kubernetes.git
cd kubernetes
git checkout tags/${KUBERNETES_VERSION}
VERSION_PRERELEASE=undistro KUBE_DOCKER_REGISTRY=docker.io/getupioundistro KUBE_RELEASE_RUN_TESTS=n make release
