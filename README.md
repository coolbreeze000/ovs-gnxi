# ovs-gnxi

A simple gNXI wrapper prototype for Open vSwitch (OVS).

Note that this implementation is not meant to be used in production, but instead only acts as a reference implementation 
to demonstrate the interaction between a gNXI Interface and OVS device.

This prototype does currently not have proper testing, nor descriptive comments.

## Run

```bash
vagrant destroy --force
vagrant up
vagrant rsync-auto
vagrant ssh ovs-gnxi
cd go/src/ovs-gnxi/scripts

./build_client.sh && docker-compose up -d --force-recreate --build client && export CONTAINER_ID_CLIENT=`docker ps -aqf 'name=client'` && docker exec -i -t $CONTAINER_ID_CLIENT bash
./gnxi_client

./build_target.sh && docker-compose up -d --force-recreate --build target && export CONTAINER_ID_TARGET=`docker ps -aqf 'name=target'` && docker exec -i -t $CONTAINER_ID_TARGET bash
tail -f /var/log/gnxi_target/gnxi_target.log
```

## Example Client Run

```bash
root@client:~# ./gnxi_client
2019-03-15 23:08:51.337971 INFO - gnxi_client.go:119 main: Started Open vSwitch gNXI client tester
2019-03-15 23:08:51.342792 INFO - gnxi_client.go:235 RunGNMICapabilitiesTests: Testing GNMI Capabilities(retrieve system capabilities)...
2019-03-15 23:08:51.405518 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/config/hostname])...
2019-03-15 23:08:51.437323 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/config/hostname]) with response value target.gnxi.lan
2019-03-15 23:08:51.437820 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/components/component[name=os]/state/description])...
2019-03-15 23:08:51.471265 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/components/component[name=os]/state/description]) with response value 2.9.2
2019-03-15 23:08:51.471585 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address])...
2019-03-15 23:08:51.500514 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address]) with response value 172.18.0.2
2019-03-15 23:08:51.500906 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port])...
2019-03-15 23:08:51.529197 INFO - gnxi_client.go:321 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port]) with response value 6653
2019-03-15 23:08:51.529649 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-03-15 23:08:51.585578 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) with response value 32
2019-03-15 23:08:51.586009 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-03-15 23:08:51.619999 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) with response value 16
2019-03-15 23:08:51.620341 INFO - gnxi_client.go:604 RunGNOIRebootTests: Testing GNOI Reboot(reboot system)...
2019-03-15 23:08:51.649648 INFO - gnxi_client.go:619 RunGNOIRebootTests: Successfully verified that target device is back up
2019-03-15 23:09:01.650165 INFO - gnxi_client.go:625 RunGNOIRebootTests: Successfully verified GNOI Reboot(reboot system)
2019-03-15 23:09:01.650252 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/config/hostname])...
2019-03-15 23:09:01.678895 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/config/hostname]) with response value target.gnxi.lan
2019-03-15 23:09:01.678960 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/components/component[name=os]/state/description])...
2019-03-15 23:09:01.712731 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/components/component[name=os]/state/description]) with response value 2.9.2
2019-03-15 23:09:01.712801 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address])...
2019-03-15 23:09:01.748954 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address]) with response value 172.18.0.2
2019-03-15 23:09:01.748997 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port])...
2019-03-15 23:09:01.779614 INFO - gnxi_client.go:321 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port]) with response value 6653
2019-03-15 23:09:01.779652 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-03-15 23:09:01.812624 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) with response value 41
2019-03-15 23:09:01.812654 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-03-15 23:09:01.850113 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) with response value 26
2019-03-15 23:09:01.850154 INFO - gnxi_client.go:634 RunGNOIGetCertificatesTests: Testing GNOI GetCertificates(get certificates)...
2019-03-15 23:09:01.884877 INFO - gnxi_client.go:648 RunGNOIGetCertificatesTests: Successfully verified GNOI GetCertificates(get certificates)
2019-03-15 23:09:01.884918 INFO - gnxi_client.go:657 RunGNOIRotateCertificatesTests: Testing GNOI RotateCertificates(rotate certificates)...
2019-03-15 23:09:07.425370 INFO - gnxi_client.go:663 RunGNOIRotateCertificatesTests: Successfully verified GNOI RotateCertificates(rotate certificates)
2019-03-15 23:09:07.425409 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/config/hostname])...
2019-03-15 23:09:07.469494 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/config/hostname]) with response value target.gnxi.lan
2019-03-15 23:09:07.469533 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/components/component[name=os]/state/description])...
2019-03-15 23:09:07.498162 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/components/component[name=os]/state/description]) with response value 2.9.2
2019-03-15 23:09:07.498190 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address])...
2019-03-15 23:09:07.530765 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address]) with response value 172.18.0.2
2019-03-15 23:09:07.530797 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port])...
2019-03-15 23:09:07.559977 INFO - gnxi_client.go:321 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port]) with response value 6653
2019-03-15 23:09:07.560061 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-03-15 23:09:07.593862 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) with response value 49
2019-03-15 23:09:07.593891 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-03-15 23:09:07.623318 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) with response value 35
2019-03-15 23:09:07.652140 INFO - gnxi_client.go:348 RunGNMISetTests: Testing GNMI Set([], [], [/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address:172.18.0.3])...
2019-03-15 23:09:07.738242 INFO - gnxi_client.go:409 RunGNMISetTests: Successfully verified GNMI Set([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address:172.18.0.3]) Update with response value 172.18.0.3
2019-03-15 23:09:07.879353 INFO - gnxi_client.go:459 RunGNMISetTests: Successfully verified Rollback GNMI Set([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address:172.18.0.3]) Update with response value 172.18.0.2
2019-03-15 23:09:07.879990 INFO - gnxi_client.go:348 RunGNMISetTests: Testing GNMI Set([], [], [/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port:6654])...
2019-03-15 23:09:07.966050 INFO - gnxi_client.go:418 RunGNMISetTests: Successfully verified GNMI Set([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port:6654]) Update with response value 6654
2019-03-15 23:09:08.092269 INFO - gnxi_client.go:468 RunGNMISetTests: Successfully verified Rollback GNMI Set([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port:6654]) Update with response value 6653
2019-03-15 23:09:08.092875 INFO - gnxi_client.go:485 RunGNMISubscribeOnceTests: Testing GNMI Subscribe ONCE([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-03-15 23:09:08.123095 INFO - gnxi_client.go:504 RunGNMISubscribeOnceTests: Successfully verified GNMI Subscribe ONCE([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) with response value 49
2019-03-15 23:09:08.123154 INFO - gnxi_client.go:485 RunGNMISubscribeOnceTests: Testing GNMI Subscribe ONCE([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-03-15 23:09:08.151625 INFO - gnxi_client.go:504 RunGNMISubscribeOnceTests: Successfully verified GNMI Subscribe ONCE([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) with response value 35
2019-03-15 23:09:08.151668 INFO - gnxi_client.go:513 RunGNMISubscribePollTests: Testing GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-03-15 23:09:08.179884 INFO - gnmi_client.go:393 SubscribePoll: Poll target again in 5 seconds
2019-03-15 23:09:08.179934 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 1/3 with response value 49
2019-03-15 23:09:13.182235 INFO - gnmi_client.go:393 SubscribePoll: Poll target again in 5 seconds
2019-03-15 23:09:13.182272 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 2/3 with response value 54
2019-03-15 23:09:18.183948 INFO - gnmi_client.go:393 SubscribePoll: Poll target again in 5 seconds
2019-03-15 23:09:18.184038 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 3/3 with response value 60
2019-03-15 23:09:18.184195 INFO - gnxi_client.go:513 RunGNMISubscribePollTests: Testing GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-03-15 23:09:18.216201 INFO - gnmi_client.go:393 SubscribePoll: Poll target again in 5 seconds
2019-03-15 23:09:18.216244 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 1/3 with response value 46
2019-03-15 23:09:23.217539 INFO - gnmi_client.go:393 SubscribePoll: Poll target again in 5 seconds
2019-03-15 23:09:23.217578 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 2/3 with response value 52
2019-03-15 23:09:28.218446 INFO - gnmi_client.go:393 SubscribePoll: Poll target again in 5 seconds
2019-03-15 23:09:28.218486 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 3/3 with response value 57
2019-03-15 23:09:28.218889 INFO - gnxi_client.go:558 RunGNMISubscribeStreamTests: Testing GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-03-15 23:09:29.381387 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 1/3 with response value 76
2019-03-15 23:09:34.384881 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 2/3 with response value 82
2019-03-15 23:09:39.391586 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 3/3 with response value 87
2019-03-15 23:09:39.391630 INFO - gnxi_client.go:558 RunGNMISubscribeStreamTests: Testing GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-03-15 23:09:49.396033 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 1/3 with response value 83
2019-03-15 23:09:54.399430 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 2/3 with response value 89
2019-03-15 23:09:59.400750 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 3/3 with response value 94
2019-03-15 23:09:59.401046 INFO - gnxi_client.go:227 main: Finished Open vSwitch gNXI client tester
2019-03-15 23:09:59.401301 INFO - gnxi_client.go:228 main: Exiting Open vSwitch gNXI client tester
```

## Example Target Run

```bash
root@target:~# tail -f /var/log/gnxi_target/gnxi_target.log
2019-03-15 23:08:30.783832 INFO - server.go:50 NewServer: Initializing gNXI Server...
2019-03-15 23:08:30.785489 INFO - gnxi_target.go:81 StartPrometheus: Starting prometheus on 0.0.0.0:8080...
2019-03-15 23:08:30.785164 INFO - manager.go:192 ActivatePackage: Cert package c5e5a1cb-8e1f-43c1-be4a-ab8e513fc667 is now active
2019-03-15 23:08:30.785712 INFO - system_broker.go:47 NewSystemBroker: Initializing OVS Client...
2019-03-15 23:08:30.785756 INFO - watchdog.go:34 NewWatchdog: Initializing Watchdog...
2019-03-15 23:08:30.785797 INFO - client.go:62 StartClient: Start OVS Client
2019-03-15 23:08:30.821959 INFO - server.go:79 createService: Start generating initial gNMI config from OVS system source...
2019-03-15 23:08:30.843191 INFO - service.go:1182 StartService: Start gNXI Service
2019-03-15 23:08:30.843310 INFO - service.go:1188 StartService: Starting to listen
2019-03-15 23:08:30.843385 INFO - service.go:1194 StartService: Starting to serve gNXI
2019-03-15 23:08:51.404443 INFO - service.go:491 Capabilities: allowed a Capabilities request
2019-03-15 23:08:51.404894 INFO - service.go:504 Capabilities: Send Capability response to client: supported_models:<name:"openconfig-interfaces" organization:"OpenConfig working group" version:"2.0.0" > supported_models:<name:"openconfig-openflow" organization:"OpenConfig working group" version:"0.1.0" > supported_models:<name:"openconfig-platform" organization:"OpenConfig working group" version:"0.5.0" > supported_models:<name:"openconfig-system" organization:"OpenConfig working group" version:"0.2.0" > supported_encodings:JSON supported_encodings:JSON_IETF gNMI_version:"0.7.0"
2019-03-15 23:08:51.436708 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:08:51.436821 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691331436818873 update:<path:<elem:<name:"system" > elem:<name:"config" > elem:<name:"hostname" > > val:<string_val:"target.gnxi.lan" > > >
2019-03-15 23:08:51.470827 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:08:51.470958 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691331470956592 update:<path:<elem:<name:"components" > elem:<name:"component" key:<key:"name" value:"os" > > elem:<name:"state" > elem:<name:"description" > > val:<string_val:"2.9.2" > > >
2019-03-15 23:08:51.499960 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:08:51.500249 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691331500246643 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.2" > > >
2019-03-15 23:08:51.528725 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:08:51.528926 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691331528924114 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6653 > > >
2019-03-15 23:08:51.585065 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:08:51.585244 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691331585242638 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:32 > > >
2019-03-15 23:08:51.619512 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:08:51.619712 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691331619710121 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:16 > > >
2019-03-15 23:08:51.648673 INFO - service.go:980 Reboot: allowed a Reboot request
2019-03-15 23:08:51.648702 INFO - service.go:990 Reboot: Send Reboot response to client:
2019-03-15 23:08:54.845504 INFO - client.go:79 StopClient: Stop OVS Client
2019-03-15 23:08:54.845542 INFO - client.go:62 StartClient: Start OVS Client
2019-03-15 23:08:54.846246 INFO - service.go:1201 StopService: Stop gNXI Service
2019-03-15 23:08:54.874118 INFO - server.go:79 createService: Start generating initial gNMI config from OVS system source...
2019-03-15 23:08:54.890406 INFO - service.go:1182 StartService: Start gNXI Service
2019-03-15 23:08:54.890531 INFO - service.go:1188 StartService: Starting to listen
2019-03-15 23:08:54.890612 INFO - service.go:1194 StartService: Starting to serve gNXI
2019-03-15 23:09:01.678503 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:01.678621 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691341678619236 update:<path:<elem:<name:"system" > elem:<name:"config" > elem:<name:"hostname" > > val:<string_val:"target.gnxi.lan" > > >
2019-03-15 23:09:01.712419 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:01.712549 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691341712548117 update:<path:<elem:<name:"components" > elem:<name:"component" key:<key:"name" value:"os" > > elem:<name:"state" > elem:<name:"description" > > val:<string_val:"2.9.2" > > >
2019-03-15 23:09:01.748438 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:01.748642 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691341748640172 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.2" > > >
2019-03-15 23:09:01.779009 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:01.779189 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691341779187430 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6653 > > >
2019-03-15 23:09:01.812103 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:01.812342 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691341812340290 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:41 > > >
2019-03-15 23:09:01.849020 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:01.849586 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691341849582644 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:26 > > >
2019-03-15 23:09:01.884393 INFO - service.go:1140 GetCertificates: allowed a GetCertificates request
2019-03-15 23:09:01.884430 INFO - service.go:1146 GetCertificates: Send GetCertificates response to client: certificate_info:<certificate_id:"c5e5a1cb-8e1f-43c1-be4a-ab8e513fc667" certificate:<type:CT_X509 certificate:"-----BEGIN CERTIFICATE-----\nMIIFRjCCAy4CCQDAPnlz9Q9eazANBgkqhkiG9w0BAQsFADBjMQswCQYDVQQGEwJB\nVDEPMA0GA1UECAwGVmllbm5hMQ8wDQYDVQQHDAZWaWVubmExDTALBgNVBAoMBFRl\nc3QxDTALBgNVBAsMBFRlc3QxFDASBgNVBAMMC2NhLmdueGkubGFuMB4XDTE5MDMx\nNTIxNDQzMVoXDTE5MDQxNDIxNDQzMVowZzELMAkGA1UEBhMCQVQxDzANBgNVBAgM\nBlZpZW5uYTEPMA0GA1UEBwwGVmllbm5hMQ0wCwYDVQQKDARUZXN0MQ0wCwYDVQQL\nDARUZXN0MRgwFgYDVQQDDA90YXJnZXQuZ254aS5sYW4wggIiMA0GCSqGSIb3DQEB\nAQUAA4ICDwAwggIKAoICAQDWeP2xafAKLNiwfHDkjhtPmuGqrhsILruLdo7QAhCr\nxC23Tum8zCKzJLEnyZMj4NlBPpoiOWlzH2g8qDK5Rf/omnBuZgfTbRUBGvWd6g1s\nR1aldjFArahPtDtI4Vpuk3fA4zwSKn/BtAS7dNw/RjZxNFs7qVcSgGks8jc21+kI\nOwhZ40nZwhJokZg7QDCJZoP/8Hl/YS4wrPdYrmUoctAicyz3IWFROsTwPdo/x49V\n0oVv9RSFuqybl/PhK4nWn+8HJbqnlRQpg9GqdKBscgNPyyu/3Eh4IP1Q+wBmjU/u\nxv+UOJVWOUS5DnXfW8V4Thkzax2HgRtgzMfs5eFaU6Cnd7oVl7NtBewgTLNoxmaS\nt6yEi9E6NsKVSE0+PshuGioLTk9g1+zrIQw3HMcS1vdixn2KRZhM2Rn/ZB9hxGzu\n+svJlcV0SLnOxCzSXVcO6l8v0bn64FuXiulZ3TGYRzAV1QC3i9HjfUR71E8wAEw9\ndnALlLLpOog3c7pLlRuyRW6gLrmwzN+isEmQwtMR7PmGgnEXWkeIaDsN1Qf0IDFC\n12oLC4rf2xmRNMkucETy0uCdQg1VysudVutm7pL/Zm+pSWQ7wocvvvZDKVahxkhw\nPm7kEKM02+dOD7Nr5PzsroTe/s/UsD5uVP3BmOj3WoRieFMAEeaQQyR812rj3Xup\nxwIDAQABMA0GCSqGSIb3DQEBCwUAA4ICAQAGsKYc42mfIAFKxVANfGLX+mIp8lTC\ncf6QJyddme0jjRwCk9Aanmt+cm7LLibzki3uoumBX3YDMZwY87hZKf1WTzJLaMAC\nY5rDcv8rd7hsj6jhMl34mU1syhMClv+zeEB+lmXXM08K3J7lLkbmnDAoNlUDwPlZ\nmHUU9NeJZclApGkA32+Ej1szZUqtDGtxOyBWUwsq0YQlC5xlkssITtyyI15xjlAV\nGdsNL2LHWB74bOQbmfE/wVdGirl5rol3ciWhwMgcwL3y2Y5b09zmO4+hs8M6kK8f\nRyYVmbZbTkAp6oFbiIGQRlOXYkkR+ttlQqbZDzNnu8QwpCcHJmUiLzT/ou9mhieh\ntK/rs9emNJkIel/Vw9VL8ktKyG157YQ0nYRtzh97CeGDLsl7zxzXDWWv7NiVVNxU\nuvfl+bQTrDtWfgY8rNOk24Hv4JYKpGqa0/nGsV4hWOqenArvLQsO5TiGG2x3lDOh\nbDld3I+HxGt9/R6UH5hr0U42UDT6wow2blwAmEHx4sbwoh7jYgIe1UStRprf8MDk\nZgSK5Rb++UNexv+H3qPUEQ7dxm3T4t/U8Opq5sfnSk61Nz5pJuHcYtLK2gHhWzwR\nugl7JEP8l44Lfw8CRRpO+waFse7qfWKdG4DYUyKMpcw8mwp6LrW2mEJJBMFIXvFe\nV5GgRkUnX9y/3A==\n-----END CERTIFICATE-----\n" > modification_time:1552691310784912540 >
2019-03-15 23:09:01.919334 INFO - service.go:1032 Rotate: allowed a Rotate request: generate_csr:<csr_params:<type:CT_X509 min_key_size:4096 key_type:KT_RSA common_name:"target.gnxi.lan" country:"AT" state:"Vienna" city:"Vienna" organization:"Test" organizational_unit:"Test" > >
2019-03-15 23:09:04.085443 INFO - manager.go:192 ActivatePackage: Cert package d7f58600-4b8e-4260-be3d-ff1641e1c8e9 is now active
2019-03-15 23:09:07.372278 INFO - client.go:79 StopClient: Stop OVS Client
2019-03-15 23:09:07.372318 INFO - client.go:62 StartClient: Start OVS Client
2019-03-15 23:09:07.425027 INFO - service.go:1140 GetCertificates: allowed a GetCertificates request
2019-03-15 23:09:07.425066 INFO - service.go:1146 GetCertificates: Send GetCertificates response to client: certificate_info:<certificate_id:"d7f58600-4b8e-4260-be3d-ff1641e1c8e9" certificate:<type:CT_X509 certificate:"-----BEGIN CERTIFICATE-----\nMIIFkzCCA3ugAwIBAgIQLbBZNCkJZV1fOhRgbD5MOjANBgkqhkiG9w0BAQsFADBj\nMQswCQYDVQQGEwJBVDEPMA0GA1UECAwGVmllbm5hMQ8wDQYDVQQHDAZWaWVubmEx\nDTALBgNVBAoMBFRlc3QxDTALBgNVBAsMBFRlc3QxFDASBgNVBAMMC2NhLmdueGku\nbGFuMB4XDTE5MDMxNTIzMDkwNFoXDTIwMDMxNDIzMDkwNFowRTELMAkGA1UEBhMC\nQVQxDTALBgNVBAoTBFRlc3QxDTALBgNVBAsTBFRlc3QxGDAWBgNVBAMTD3Rhcmdl\ndC5nbnhpLmxhbjCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAM5v1Ssh\nam9jLV3xP6tIYPYseQb1ORQPY6wKIofdefiovyJT+0StJxsjizi+K+61wJAh2oN3\nGWyaTImfoniG3Vog9ppLhdTqwtAMrm3Phm8aK7ZrnYeblRYa96Nisue0YiRLHbeG\n2avYoe1wq3j8MYHON2M9pxeZ5AheHp2BmyU//h/uDroDjv9smqtXYpzhCDBh2Gk5\nWV4uB/AJXeZGwWZava8gd50++Bsp1jtj6pT9tHd5PK3jRJIdILsstVmHAeAnUhxx\n6usd3bxWyJ4ULv6dwDEgbgO0jvBzIVKtHDNYHzh91ikofDzm49W8W7TQVcGfupzk\n0/jXTGQ4/6x1jBHn3Sx2jE09lAoZryKYHgoMe3tE14m06SAVU0u8fm+OGi45vG3E\nR+JA0XsGbV34VcDpVea+xieqTLHSDLAgDwXsBvlyA76ADcDn+Vuw7DCqMat7g3bO\niJSNIcxG9AEVtuw1nq31SLHgivEfcJbPFn7VNJ8LYO5FgPFb0BCxR0KncKF5pMjs\nGVKCUdFOtv81eLW98OqXW5xwEIL/bBk1RZJLagCUDT9XrLTbMGOAo3u2gVwGL8P4\nA70Ecst0GYZANSr8t2g3aHImeR3/bxbNU7Tiru6qTN6wpoX0TXSss5nYDBaQa10K\nrQVDz4xu3LYKme4zm8lIb+rm8/Wz17z6FoqJAgMBAAGjYTBfMA4GA1UdDwEB/wQE\nAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw\nAwIBBTAdBgNVHQ4EFgQULCaGVWvV3ZfPZ99o5v7iX851a10wDQYJKoZIhvcNAQEL\nBQADggIBAGNLpaigsSsObjyvRXygszpOhh0MfIqcSb15VNhNs/voMVNUP4bCwxzr\nuXVZ1YQddSE8f5E8VrQKfbtaDSueqX29Ak+ZDTOx9iPKUG5uBiUxJu75T0mnNEAi\nD6RUy+9DMwTeaZKNMCPnphZZz78bUyaHnznjB5xIRqTWemBA6XB4kIVJIwL1+gfa\nQbHDzI+dzMb9JevkfBv1ApVBBNDM8F/Mw5OiIlBcvipn0Qr8rvvuMMlOa1ELEX5j\nzzuczfhx7N7b4PS0KqOpS61MuIKEAcF3nN8Jlq3yAF1NFTJtH8a5X6lu8Z0wTSyg\nvtmovDcDibcL6LoYicKRYwp5D/xlDRbPRSShdyWzrm+4I9whclmdlPqjBBeTHs/u\nd9WDnzxpuIax1cPxLcSjM6lbdzrdfhMI99Vk+hZ62yMREwOszPyQLJ63FAgbxi1S\nnAOExJIE8rE2AoLrlTq2W2OLyv0XC4WO1taVvZPbEqe1JNMUdv9DONqLlrFAO9DS\nuvmIgebWOI8L88DmqllSvIjgKwO088wCXadJrzMdmUjvbI9LopsKsGUKYLl1elvD\nBwcHfXeLTgTuafjXawREw830YWNUF2nosHLpzpzKpgL6D+GoHrOJb1GR+1Xj+zaD\nShRmghMcprGt/LzruQPqvsCnUVHVDZT2/5xWoDFpSHOQuBgi8TXm\n-----END CERTIFICATE-----\n" > modification_time:1552691344083782681 >
2019-03-15 23:09:07.468877 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:07.469032 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691347469030410 update:<path:<elem:<name:"system" > elem:<name:"config" > elem:<name:"hostname" > > val:<string_val:"target.gnxi.lan" > > >
2019-03-15 23:09:07.497740 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:07.497920 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691347497917842 update:<path:<elem:<name:"components" > elem:<name:"component" key:<key:"name" value:"os" > > elem:<name:"state" > elem:<name:"description" > > val:<string_val:"2.9.2" > > >
2019-03-15 23:09:07.530320 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:07.530517 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691347530515274 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.2" > > >
2019-03-15 23:09:07.559583 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:07.559764 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691347559762081 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6653 > > >
2019-03-15 23:09:07.593374 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:07.593594 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691347593592470 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:49 > > >
2019-03-15 23:09:07.622735 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:07.623065 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691347623062285 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:35 > > >
2019-03-15 23:09:07.680397 INFO - service.go:639 Set: allowed a Set request
2019-03-15 23:09:07.689995 INFO - client.go:212 SyncChangesToRemote: target is in inconsistent state with OVS device, syncing Controller
2019-03-15 23:09:07.698356 INFO - service.go:694 Set: Send Set response to client: response:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > op:UPDATE >
2019-03-15 23:09:07.737813 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:07.737975 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691347737973765 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.3" > > >
2019-03-15 23:09:07.766946 INFO - service.go:639 Set: allowed a Set request
2019-03-15 23:09:07.774190 INFO - client.go:212 SyncChangesToRemote: target is in inconsistent state with OVS device, syncing Controller
2019-03-15 23:09:07.831508 INFO - service.go:694 Set: Send Set response to client: response:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > op:UPDATE >
2019-03-15 23:09:07.878935 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:07.879117 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691347879115429 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.2" > > >
2019-03-15 23:09:07.908775 INFO - service.go:639 Set: allowed a Set request
2019-03-15 23:09:07.917125 INFO - client.go:212 SyncChangesToRemote: target is in inconsistent state with OVS device, syncing Controller
2019-03-15 23:09:07.928626 INFO - service.go:694 Set: Send Set response to client: response:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > op:UPDATE >
2019-03-15 23:09:07.965626 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:07.965792 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691347965790386 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6654 > > >
2019-03-15 23:09:07.995106 INFO - service.go:639 Set: allowed a Set request
2019-03-15 23:09:08.000819 INFO - client.go:212 SyncChangesToRemote: target is in inconsistent state with OVS device, syncing Controller
2019-03-15 23:09:08.023291 INFO - service.go:694 Set: Send Set response to client: response:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > op:UPDATE >
2019-03-15 23:09:08.091750 INFO - service.go:516 Get: allowed a Get request
2019-03-15 23:09:08.092013 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1552691348092011127 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6653 > > >
2019-03-15 23:09:08.122326 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > > mode:ONCE encoding:JSON_IETF >
2019-03-15 23:09:08.122445 INFO - service.go:877 subscribeOnce: serving subscribe ONCE
2019-03-15 23:09:08.122723 INFO - service.go:885 subscribeOnce: Send Subscribe ONCE response to client: update:<timestamp:1552691348122718202 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:49 > > >
2019-03-15 23:09:08.151103 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > > mode:ONCE encoding:JSON_IETF >
2019-03-15 23:09:08.151173 INFO - service.go:877 subscribeOnce: serving subscribe ONCE
2019-03-15 23:09:08.151382 INFO - service.go:885 subscribeOnce: Send Subscribe ONCE response to client: update:<timestamp:1552691348151378131 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:35 > > >
2019-03-15 23:09:08.179373 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > > mode:POLL encoding:JSON_IETF >
2019-03-15 23:09:08.179490 INFO - service.go:897 subscribePoll: serving subscribe POLL
2019-03-15 23:09:08.179682 INFO - service.go:903 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1552691348179679068 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:49 > > >
2019-03-15 23:09:13.181792 INFO - service.go:940 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1552691353181781620 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:54 > > >
2019-03-15 23:09:18.183145 INFO - service.go:940 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1552691358182962563 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:60 > > >
2019-03-15 23:09:18.215045 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > > mode:POLL encoding:JSON_IETF >
2019-03-15 23:09:18.215195 INFO - service.go:897 subscribePoll: serving subscribe POLL
2019-03-15 23:09:18.215436 INFO - service.go:903 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1552691358215419776 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:46 > > >
2019-03-15 23:09:23.217325 INFO - service.go:940 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1552691363217316778 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:52 > > >
2019-03-15 23:09:28.218248 INFO - service.go:940 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1552691368218240121 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:57 > > >
2019-03-15 23:09:28.254051 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > > encoding:JSON_IETF >
2019-03-15 23:09:28.254122 INFO - service.go:952 subscribeStream: serving subscribe STREAM
2019-03-15 23:09:29.380996 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1552691369380993415 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:76 > > >
2019-03-15 23:09:34.384525 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1552691374384520449 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:82 > > >
2019-03-15 23:09:39.391151 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1552691379391148133 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:87 > > >
2019-03-15 23:09:39.421225 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > > encoding:JSON_IETF >
2019-03-15 23:09:39.421316 INFO - service.go:952 subscribeStream: serving subscribe STREAM
2019-03-15 23:09:44.388292 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1552691384388281894 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:92 > > >
2019-03-15 23:09:49.395363 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1552691389395359029 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:83 > > >
2019-03-15 23:09:54.398814 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1552691394398811109 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:89 > > >
2019-03-15 23:09:59.400363 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1552691399400360777 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:94 > > >
2019-03-15 23:10:04.444316 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1552691404444312967 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:99 > > >
```

## Generate Go Bindings for Open vSwitch
```bash
go run vendor/github.com/openconfig/ygot/generator/generator.go -path=yang -output_file=oc.go -package_name=openvswitch -generate_fakeroot -fakeroot_name=device -compress_paths=true -exclude_modules=ietf-interfaces yang/openconfig-interfaces.yang
```