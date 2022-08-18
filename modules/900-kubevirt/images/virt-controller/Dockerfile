# Source https://github.com/kubevirt/kubevirt/blob/main/hack/dockerized#L15
ARG KUBEVIRT_BUILDER_IMAGE="quay.io/kubevirt/builder:2206141426-2d31aa580"

FROM $KUBEVIRT_BUILDER_IMAGE as builder
COPY unpack-bundle.sh unpack-bundle.sh

ARG KUBEVIRT_VERSION=0.55.0

# Copy patches
COPY ./patches /patches

RUN git clone --branch v${KUBEVIRT_VERSION} https://github.com/kubevirt/kubevirt /kubevirt
RUN cd /kubevirt \
 && mkdir -p _out \
 && git apply /patches/*.patch \
 && make bazel-build-image-bundle KUBEVIRT_RUN_UNNESTED=true
RUN tar --one-top-level -xf /kubevirt/_out/virt-components-bundle.tar
RUN mkdir images \
 && cd images \
 && /unpack-bundle.sh /virt-components-bundle/

# ------------------------------------------------------------------------------

FROM scratch
COPY --from=builder /images/kubevirt/virt-controller:latest /
USER 1001
ENTRYPOINT ["/usr/bin/virt-controller"]
