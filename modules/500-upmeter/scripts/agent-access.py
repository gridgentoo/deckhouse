#!/usr/bin/env python3

import subprocess
from typing import List


class CanI(object):
    def __init__(self, verb: str, resource: str, subresource: str = None, name: str = None,
                 namespace: str = None) -> None:
        super().__init__()
        self.service_account = "system:serviceaccount:d8-upmeter:upmeter-agent"
        self.verb = verb
        self.resource = resource
        self.subresource = subresource
        self.name = name
        self.namespace = namespace

    def command(self):
        """
        kubectl auth can-i --help

        Check whether an action is allowed.

        VERB is a logical Kubernetes API verb like 'get', 'list', 'watch', 'delete', etc. TYPE is a
        Kubernetes resource.

        Shortcuts and groups will be resolved. NONRESOURCEURL is a partial URL starts with "/". NAME is
        the name of a particular

        Kubernetes resource.

        Examples:
            # Check to see if I can create pods in any namespace
            kubectl auth can-i create pods --all-namespaces

            # Check to see if I can list deployments in my current namespace
            kubectl auth can-i list deployments.apps

            # Check to see if I can do everything in my current namespace ("*" means all)
            kubectl auth can-i '*' '*'

            # Check to see if I can get the job named "bar" in namespace "foo"
            kubectl auth can-i list jobs.batch/bar -n foo

            # Check to see if I can read pod logs
            kubectl auth can-i get pods --subresource=log

            # Check to see if I can access the URL /logs/
            kubectl auth can-i get /logs/

            # List all allowed actions in namespace "foo"
            kubectl auth can-i --list --namespace=foo

        Options:
        -A, --all-namespaces=false: If true, check the specified action in all namespaces.
            --list=false: If true, prints all allowed actions.
            --no-headers=false: If true, prints allowed actions without headers
        -q, --quiet=false: If true, suppress output and just return the exit code.
            --subresource='': SubResource such as pod/log or deployment/scale

        Usage:
        kubectl auth can-i VERB [TYPE | TYPE/NAME | NONRESOURCEURL] [options]

        Use "kubectl options" for a list of global command-line options (applies to all commands).
        """
        cmd = [
            "kubectl", "auth", "can-i",
            "--as={}".format(self.service_account),
            self.verb
        ]

        if self.name is not None:
            cmd.append("{}/{}".format(self.resource, self.name))
        else:
            cmd.append(self.resource)

        if self.subresource is not None:
            cmd.append("--subresource")
            cmd.append(self.subresource)

        if self.namespace is None:
            # Omit namespace warnings, see https://github.com/kubernetes/kubernetes/pull/76014#issuecomment-481244715
            cmd.append("--all-namespaces")
        else:
            cmd.append("-n")
            cmd.append(self.namespace)

        return cmd

    def claim(self):
        claim = "{} {}".format(self.verb, self.resource)

        if self.subresource is not None:
            claim += "/{}".format(self.subresource)

        if self.name is not None:
            claim += "/{}".format(self.name)

        if self.namespace is not None:
            claim += " in {}".format(self.namespace)

        return claim


def call(cmd: List[str]) -> bytes:
    """
    :cmd: is a list to run a command

    :return: "yes" or "no"
    """
    p = subprocess.run(cmd, stdout=subprocess.PIPE, universal_newlines=True)
    return p.stdout.strip()


def can_i(verb, resource, subresource=None, name=None, namespace=None):
    c = CanI(verb, resource, subresource=subresource, name=name, namespace=namespace)
    claim = c.claim()
    cmd = c.command()

    # color_grey = subprocess.check_output("tput setaf 8".split())
    # color_reset = subprocess.check_output("tput sgr0".split())
    print(" ".join(cmd))

    # assert call(cmd) == "yes", "Cannot {}".format(claim)

    if call(cmd) != "yes":
        print("Cannot {}".format(claim))
        return

    print("Can {}".format(claim))


if __name__ == '__main__':
    # Listing and deletion is usually required for garbage collection

    #
    # Control plane
    #

    # control-plane/access
    can_i("get", "/version")

    # control-plane/basic-functionality
    can_i("list", "configmaps", namespace="d8-upmeter")
    can_i("create", "configmaps", namespace="d8-upmeter")
    can_i("get", "configmaps", namespace="d8-upmeter")
    can_i("delete", "configmaps", namespace="d8-upmeter")

    # control-plane/namespace
    can_i("create", "namespaces")
    can_i("list", "namespaces")
    can_i("delete", "namespaces")

    # control-plane/controller-manager
    can_i("list", "deployments", namespace="d8-upmeter")
    can_i("create", "deployments", namespace="d8-upmeter")
    can_i("get", "deployments", namespace="d8-upmeter")
    can_i("delete", "deployments", namespace="d8-upmeter")

    # control-plane/scheduler
    can_i("list", "pods", namespace="d8-upmeter")
    can_i("create", "pods", namespace="d8-upmeter")
    can_i("get", "pods", namespace="d8-upmeter")
    can_i("delete", "pods", namespace="d8-upmeter")

    #
    # Deckhouse
    #

    # deckhouse/cluster-configuration
    can_i("list", "pods", namespace="d8-system")
    can_i("create", "upmeterhookprobes")
    can_i("get", "upmeterhookprobes")
    can_i("update", "upmeterhookprobes")
    can_i("watch", "upmeterhookprobes")

    #
    # Load balancing
    #

    # load-balancing/load-balancer-configuration
    can_i("list", "pods", namespace="d8-cloud-provider-aws")
    can_i("list", "pods", namespace="d8-cloud-provider-azure")
    can_i("list", "pods", namespace="d8-cloud-provider-gcp")
    can_i("list", "pods", namespace="d8-cloud-provider-openstack")
    can_i("list", "pods", namespace="d8-cloud-provider-vsphere")
    can_i("list", "pods", namespace="d8-cloud-provider-yandex")

    # load-balancing/metallb
    can_i("list", "pods", namespace="d8-metallb")

    #
    # Monitoring and autoscaling
    #

    # monitoring-and-autoscaling/*
    can_i("list", "pods", namespace="d8-monitoring")

    # monitoring-and-autoscaling/prometheus
    # monitoring-and-autoscaling/key-metrics-present
    can_i("get", "prometheuses", name="main", subresource="http", namespace="d8-monitoring")

    # monitoring-and-autoscaling/trickster
    can_i("get", "deploy", name="trickster", subresource="http", namespace="d8-monitoring")

    # monitoring-and-autoscaling/prometheus-metrics-adapter
    can_i("get", "/apis/custom.metrics.k8s.io/v1beta1/namespaces/d8-upmeter/metrics/memory_1m")

    # monitoring-and-autoscaling/vertical-pod-autoscaler
    can_i("list", "pods", namespace="kube-system")

    # monitoring-and-autoscaling/metrics-sources
    can_i("get", "daemonsets", namespace="d8-monitoring")

    #
    # Scaling
    #

    # scaling/cluster-scaling
    # scaling/cluster-autoscaler
    can_i("list", "pods", namespace="d8-cloud-instance-manager")
    can_i("list", "pods", namespace="d8-cloud-provider-aws")
    can_i("list", "pods", namespace="d8-cloud-provider-azure")
    can_i("list", "pods", namespace="d8-cloud-provider-gcp")
    can_i("list", "pods", namespace="d8-cloud-provider-openstack")
    can_i("list", "pods", namespace="d8-cloud-provider-vsphere")
    can_i("list", "pods", namespace="d8-cloud-provider-yandex")

    #
    # Synthetic
    #

    can_i("get", "statefulsets", subresource="http", name="smoke-mini-a", namespace="d8-upmeter")
    can_i("get", "statefulsets", subresource="http", name="smoke-mini-b", namespace="d8-upmeter")
    can_i("get", "statefulsets", subresource="http", name="smoke-mini-c", namespace="d8-upmeter")
    can_i("get", "statefulsets", subresource="http", name="smoke-mini-d", namespace="d8-upmeter")
    can_i("get", "statefulsets", subresource="http", name="smoke-mini-e", namespace="d8-upmeter")
