#!/usr/bin/env sh

# CALL THIS SCRIPT FROM THE MODULE DIRECTORY
#  ./scripts/get_manifests.sh

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

ARGOCD_REPO="${HOME}/dev/flant/argoproj/argo-cd"
MANIFESTS="${ARGOCD_REPO}/manifests/install.yaml"
HA_MANIFESTS="${ARGOCD_REPO}/manifests/ha/install.yaml"

pull_manifests() {
        MANIFESTS=$1

        yq eval-all 'select(.kind == "CustomResourceDefinition") | .' $MANIFESTS |
                yq e --no-doc -s '"crd-" + .spec.names.singular' -

        yq eval-all 'select(.kind != "CustomResourceDefinition") | .' $MANIFESTS |
                yq e --no-doc -s '.metadata.name' -

        # move CRD

        mv crd-*.yml crds &&
                pushd crds &&
                rename 's/^crd-//g' * &&
                popd

        # move other manifests
        mkdir -p templates/argocd-application-controller
        mv argocd-application-controller*.yml templates/argocd-application-controller
        mv argocd-metrics.yml templates/argocd-application-controller

        mkdir -p templates/argocd-repo-server
        mv argocd-repo-server*.yml templates/argocd-repo-server

        mkdir -p templates/argocd-server
        mv argocd-server*.yml templates/argocd-server

        mkdir -p templates/argocd-dex
        mv argocd-dex*.yml templates/argocd-dex
        rm -rf templates/argocd-dex # We use or own dex

        mkdir -p templates/redis
        mv argocd-redis*.yml templates/redis
        pushd templates/redis && rename 's/^(.*)$/ha-$1/g' *-ha* && rename 's/-ha//' *-ha* && popd

        # all other manifests
        mv argocd-*.yml templates/
}

# clean existing manifests
mkdir -p crds
mkdir -p templates
rm -rf templates/* crds/*

pull_manifests "${MANIFESTS}"
# pull_manifests "${HA_MANIFESTS}"
