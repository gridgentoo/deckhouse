---
title: "Модуль local-path-provisioner"
---

Local Path Provisioner позволяет пользователям Kubernetes использовать локальное хранилище на узлах. 

## Как это работает ?
Для каждого CR ```LocalPathProvisioner``` создается соответствующий ```StorageClass```.
Допустимая топология для SC вычисляется на основе списка нодгрупп из CR.
Топология используется при шедулинге подов.
Когда под заказывает диск, то создается ```HostPath``` PV, а ```Provisioner``` создает на нужном узле локальную папку по пути, состоящем
из параметра ```path``` CR, имени PV и имени PVC
(например, ```/opt/local-path-provisioner/pvc-d9bd3878-f710-417b-a4b3-38811aa8aac1_d8-monitoring_prometheus-main-db-prometheus-main-0```).

## Ограничения
Ограничение размера диска не поддерживается для локальных томов.