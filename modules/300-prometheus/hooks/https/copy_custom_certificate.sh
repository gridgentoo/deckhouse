#!/bin/bash

source /deckhouse/shell_lib.sh

function __config__() {
  common_hooks::https::copy_custom_certificate::config
}

function __main__() {
  common_hooks::https::copy_custom_certificate::main "kube-prometheus" "ingress-tls"
}

hook::run "$@"
