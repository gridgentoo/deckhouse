# Changelog v1.36

## Know before update


 - All ingress nginx controllers with not-specified version (0.33) will restart and upgrade to 1.1
 - All of the controllers will be restarted.

## Features


 - **[candi]** Set `maxAllowed` and `minAllowed` to all VPA objects. Set resources requests for all controllers if VPA is off.   Added `global.modules.resourcesRequests.controlPlane` values. `global.modules.resourcesRequests.EveryNode` and `global.modules.resourcesRequests.masterNode` values are deprecated. [#1918](https://github.com/deckhouse/deckhouse/pull/1918)
    All of the controllers will be restarted.
 - **[ingress-nginx]** Change default ingress nginx controller version to 1.1 [#2267](https://github.com/deckhouse/deckhouse/pull/2267)
    All ingress nginx controllers with not-specified version (0.33) will restart and upgrade to 1.1
 - **[log-shipper]** Refactor transforms composition, improve efficiency and fix destination transforms. [#2050](https://github.com/deckhouse/deckhouse/pull/2050)

## Fixes


 - **[node-manager]** Change cluster autoscaler timeouts to avoid node flapping [#2279](https://github.com/deckhouse/deckhouse/pull/2279)

