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

ARGOCD_REPO="${HOME}/dev/flant/argoproj/argo-cd"
MANIFESTS="${ARGOCD_REPO}/manifests/ha/install.yaml"
MANIFESTS="${ARGOCD_REPO}/manifests/install.yaml"

yq eval-all 'select(.kind == "CustomResourceDefinition") | .' $MANIFESTS |
        yq e --no-doc -s '"crd-" + .spec.names.singular' -

yq eval-all 'select(.kind != "CustomResourceDefinition") | .' $MANIFESTS |
        yq e --no-doc -s '.metadata.name' -

# clean existing manifests
rm -rf templates/* crds/*

# move CRD
mkdir -p crds &&
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
pushd templates/redis && rename 's/-ha//g' * && popd

# all other manifests
mv argocd-*.yml templates/
