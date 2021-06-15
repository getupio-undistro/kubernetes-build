#!/usr/bin/env bash

KUBERNETES_VERSION=$1

git clone https://github.com/kubernetes/kubernetes.git
cd kubernetes
git checkout tags/${KUBERNETES_VERSION}
VERSION_PRERELEASE=undistro KUBE_RELEASE_RUN_TESTS=n make release
