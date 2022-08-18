#!/bin/sh
# Copyright 2021 Flant JSC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

##########################################################################
# CALL THIS SCRIPT FROM THE MODULE DIRECTORY
#  ./scripts/get_manifests.sh
#
# Dependencies:
#   yq                  https://github.com/mikefarah/yq
#
#   argocd repo         https://github.com/argoproj/argo-cd
#                      (checkout specific tag)
#
#   rename             famous perl script
#

set -euo pipefail
set -x

VERSION="v2.4.9"

# First, clone the repo
ARGOCD_REPO="${HOME}/dev/flant/argoproj/argo-cd"
MANIFESTS="${ARGOCD_REPO}/manifests/install.yaml"
HA_MANIFESTS="${ARGOCD_REPO}/manifests/ha/install.yaml"

mkdir -p "${ARGOCD_REPO}"
git clone git@github.com:argoproj/argo-cd.git "${ARGOCD_REPO}" || true
pushd $ARGOCD_REPO &&
  git clean -df &&
  git reset --hard &&
  git fetch --all --prune &&
  git checkout $VERSION &&
  popd

CRD_ROOT=crds
MANIFESTS_ROOT=templates/argocd

pull_manifests() {
  MANIFESTS=$1

  yq eval-all 'select(.kind == "CustomResourceDefinition") | .' $MANIFESTS |
    yq e --no-doc -s '"crd-" + .spec.names.singular' -

  yq eval-all 'select(.kind != "CustomResourceDefinition") | .' $MANIFESTS |
    yq e --no-doc -s '.metadata.name' -

  # move CRD
  mv crd-*.yml ${CRD_ROOT} &&
    pushd ${CRD_ROOT} &&
    rename 's/^crd-(.*)/argocd-$1/g' * &&
    popd

  # move other manifests
  mkdir -p ${MANIFESTS_ROOT}/argocd-application-controller
  mv argocd-application-controller*.yml ${MANIFESTS_ROOT}/argocd-application-controller
  mv argocd-applicationset-controller*.yml ${MANIFESTS_ROOT}/argocd-application-controller
  mv argocd-metrics.yml ${MANIFESTS_ROOT}/argocd-application-controller

  mkdir -p ${MANIFESTS_ROOT}/argocd-notifications
  mv argocd-notifications*.yml ${MANIFESTS_ROOT}/argocd-notifications

  mkdir -p ${MANIFESTS_ROOT}/argocd-repo-server
  mv argocd-repo-server*.yml ${MANIFESTS_ROOT}/argocd-repo-server

  mkdir -p ${MANIFESTS_ROOT}/argocd-server
  mv argocd-server*.yml ${MANIFESTS_ROOT}/argocd-server

  mkdir -p ${MANIFESTS_ROOT}/argocd-dex
  mv argocd-dex*.yml ${MANIFESTS_ROOT}/argocd-dex
  # We use our own dex
  rm -rf ${MANIFESTS_ROOT}/argocd-dex

  mkdir -p ${MANIFESTS_ROOT}/redis
  mv argocd-redis*.yml ${MANIFESTS_ROOT}/redis
  pushd ${MANIFESTS_ROOT}/redis && rename 's/^(.*)$/ha-$1/g' *-ha* && rename 's/-ha//' *-ha* && popd

  # all other manifests
  mv argocd-*.yml ${MANIFESTS_ROOT}/
}

# clean existing manifests
mkdir -p $CRD_ROOT
mkdir -p $MANIFESTS_ROOT
rm -rf ${MANIFESTS_ROOT}/* crds/argocd-*

pull_manifests "${MANIFESTS}"
# pull_manifests "${HA_MANIFESTS}"
