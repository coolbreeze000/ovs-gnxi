# ovs-gnxi

A simple gNXI wrapper prototype for Open vSwitch (OVS).

Note that this implementation is not meant to be used in production, but instead only acts as a reference implementation 
to demonstrate the interaction between a gNXI Interface and OVS device.

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
2019-02-19 22:42:38.190080 INFO - gnxi_client.go:119 main: Started Open vSwitch gNXI client tester
2019-02-19 22:42:38.194631 INFO - gnxi_client.go:235 RunGNMICapabilitiesTests: Testing GNMI Capabilities(retrieve system capabilities)...
2019-02-19 22:42:38.272079 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/config/hostname])...
2019-02-19 22:42:38.300722 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/config/hostname]) with response value target.gnxi.lan
2019-02-19 22:42:38.300755 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/components/component[name=os]/state/description])...
2019-02-19 22:42:38.330107 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/components/component[name=os]/state/description]) with response value 2.9.2
2019-02-19 22:42:38.330525 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address])...
2019-02-19 22:42:38.359707 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address]) with response value 172.18.0.2
2019-02-19 22:42:38.359748 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port])...
2019-02-19 22:42:38.390312 INFO - gnxi_client.go:321 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port]) with response value 6653
2019-02-19 22:42:38.390740 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-02-19 22:42:38.419304 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) with response value 21
2019-02-19 22:42:38.420199 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-02-19 22:42:38.448273 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) with response value 16
2019-02-19 22:42:38.448366 INFO - gnxi_client.go:604 RunGNOIRebootTests: Testing GNOI Reboot(reboot system)...
2019-02-19 22:42:38.484736 INFO - gnxi_client.go:619 RunGNOIRebootTests: Successfully verified that target device is back up
2019-02-19 22:42:48.485225 INFO - gnxi_client.go:625 RunGNOIRebootTests: Successfully verified GNOI Reboot(reboot system)
2019-02-19 22:42:48.485349 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/config/hostname])...
2019-02-19 22:42:48.519461 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/config/hostname]) with response value target.gnxi.lan
2019-02-19 22:42:48.519490 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/components/component[name=os]/state/description])...
2019-02-19 22:42:48.551019 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/components/component[name=os]/state/description]) with response value 2.9.2
2019-02-19 22:42:48.551526 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address])...
2019-02-19 22:42:48.580352 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address]) with response value 172.18.0.2
2019-02-19 22:42:48.580610 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port])...
2019-02-19 22:42:48.614810 INFO - gnxi_client.go:321 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port]) with response value 6653
2019-02-19 22:42:48.615283 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-02-19 22:42:48.652022 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) with response value 29
2019-02-19 22:42:48.652289 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-02-19 22:42:48.689544 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) with response value 29
2019-02-19 22:42:48.689922 INFO - gnxi_client.go:634 RunGNOIGetCertificatesTests: Testing GNOI GetCertificates(get certificates)...
2019-02-19 22:42:48.730317 INFO - gnxi_client.go:648 RunGNOIGetCertificatesTests: Successfully verified GNOI GetCertificates(get certificates)
2019-02-19 22:42:48.730820 INFO - gnxi_client.go:657 RunGNOIRotateCertificatesTests: Testing GNOI Reboot(rotate certificates)...
2019-02-19 22:42:54.746294 INFO - gnxi_client.go:663 RunGNOIRotateCertificatesTests: Successfully verified GNOI Rotate Certificates(rotate certificates)
2019-02-19 22:42:54.746369 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/config/hostname])...
2019-02-19 22:42:54.795363 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/config/hostname]) with response value target.gnxi.lan
2019-02-19 22:42:54.795670 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/components/component[name=os]/state/description])...
2019-02-19 22:42:54.827298 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/components/component[name=os]/state/description]) with response value 2.9.2
2019-02-19 22:42:54.827401 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address])...
2019-02-19 22:42:54.873111 INFO - gnxi_client.go:312 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address]) with response value 172.18.0.2
2019-02-19 22:42:54.873229 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port])...
2019-02-19 22:42:54.905674 INFO - gnxi_client.go:321 RunGNMIGetTests: Successfully verified GNMI Get([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port]) with response value 6653
2019-02-19 22:42:54.905720 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-02-19 22:42:54.935558 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) with response value 38
2019-02-19 22:42:54.935617 INFO - gnxi_client.go:298 RunGNMIGetTests: Testing GNMI Get([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-02-19 22:42:54.965116 INFO - gnxi_client.go:327 RunGNMIGetTests: Successfully verified GNMI Subscribe([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) with response value 34
2019-02-19 22:42:54.995411 INFO - gnxi_client.go:348 RunGNMISetTests: Testing GNMI Set([], [], [/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address:172.18.0.3])...
2019-02-19 22:42:55.075456 INFO - gnxi_client.go:409 RunGNMISetTests: Successfully verified GNMI Set([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address:172.18.0.3]) Update with response value 172.18.0.3
2019-02-19 22:42:55.170692 INFO - gnxi_client.go:459 RunGNMISetTests: Successfully verified Rollback GNMI Set([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address:172.18.0.3]) Update with response value 172.18.0.2
2019-02-19 22:42:55.171390 INFO - gnxi_client.go:348 RunGNMISetTests: Testing GNMI Set([], [], [/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port:6654])...
2019-02-19 22:42:55.250397 INFO - gnxi_client.go:418 RunGNMISetTests: Successfully verified GNMI Set([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port:6654]) Update with response value 6654
2019-02-19 22:42:55.345064 INFO - gnxi_client.go:468 RunGNMISetTests: Successfully verified Rollback GNMI Set([/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/port:6654]) Update with response value 6653
2019-02-19 22:42:55.345639 INFO - gnxi_client.go:485 RunGNMISubscribeOnceTests: Testing GNMI Subscribe ONCE([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-02-19 22:42:55.375560 INFO - gnxi_client.go:504 RunGNMISubscribeOnceTests: Successfully verified GNMI Subscribe ONCE([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) with response value 38
2019-02-19 22:42:55.375614 INFO - gnxi_client.go:485 RunGNMISubscribeOnceTests: Testing GNMI Subscribe ONCE([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-02-19 22:42:55.404432 INFO - gnxi_client.go:504 RunGNMISubscribeOnceTests: Successfully verified GNMI Subscribe ONCE([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) with response value 34
2019-02-19 22:42:55.404778 INFO - gnxi_client.go:513 RunGNMISubscribePollTests: Testing GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-02-19 22:42:55.435374 INFO - gnmi_client.go:378 SubscribePoll: Poll target again in 5 seconds
2019-02-19 22:42:55.435432 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 1/3 with response value 38
2019-02-19 22:43:00.437127 INFO - gnmi_client.go:378 SubscribePoll: Poll target again in 5 seconds
2019-02-19 22:43:00.437168 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 2/3 with response value 43
2019-02-19 22:43:05.438953 INFO - gnmi_client.go:378 SubscribePoll: Poll target again in 5 seconds
2019-02-19 22:43:05.439031 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 3/3 with response value 49
2019-02-19 22:43:05.439053 INFO - gnxi_client.go:513 RunGNMISubscribePollTests: Testing GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-02-19 22:43:05.468335 INFO - gnmi_client.go:378 SubscribePoll: Poll target again in 5 seconds
2019-02-19 22:43:05.468374 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 1/3 with response value 45
2019-02-19 22:43:10.470136 INFO - gnmi_client.go:378 SubscribePoll: Poll target again in 5 seconds
2019-02-19 22:43:10.470189 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 2/3 with response value 51
2019-02-19 22:43:15.471577 INFO - gnmi_client.go:378 SubscribePoll: Poll target again in 5 seconds
2019-02-19 22:43:15.471652 INFO - gnxi_client.go:537 RunGNMISubscribePollTests: Successfully verified GNMI Subscribe POLL([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 3/3 with response value 57
2019-02-19 22:43:15.471682 INFO - gnxi_client.go:558 RunGNMISubscribeStreamTests: Testing GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts])...
2019-02-19 22:43:16.712540 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 1/3 with response value 66
2019-02-19 22:43:21.714405 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 2/3 with response value 71
2019-02-19 22:43:26.718851 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/in-pkts]) 3/3 with response value 75
2019-02-19 22:43:26.718906 INFO - gnxi_client.go:558 RunGNMISubscribeStreamTests: Testing GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts])...
2019-02-19 22:43:36.719874 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 1/3 with response value 84
2019-02-19 22:43:41.724638 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 2/3 with response value 89
2019-02-19 22:43:46.724962 INFO - gnxi_client.go:582 RunGNMISubscribeStreamTests: Successfully verified GNMI Subscribe STREAM([/interfaces/interface[name=sw1-eth1]/state/counters/out-pkts]) 3/3 with response value 94
2019-02-19 22:43:46.725010 INFO - gnxi_client.go:227 main: Finished Open vSwitch gNXI client tester
2019-02-19 22:43:46.725030 INFO - gnxi_client.go:228 main: Exiting Open vSwitch gNXI client tester
```

## Example Target Run

```bash
root@target:~# tail -f /var/log/gnxi_target/gnxi_target.log
2019-02-19 22:42:29.791287 INFO - server.go:50 NewServer: Initializing gNXI Server...
2019-02-19 22:42:29.791386 INFO - gnxi_target.go:81 StartPrometheus: Starting prometheus on 0.0.0.0:8080...
2019-02-19 22:42:29.795800 INFO - manager.go:192 ActivatePackage: Cert package c5e5a1cb-8e1f-43c1-be4a-ab8e513fc667 is now active
2019-02-19 22:42:29.795838 INFO - system_broker.go:47 NewSystemBroker: Initializing OVS Client...
2019-02-19 22:42:29.795897 INFO - watchdog.go:34 NewWatchdog: Initializing Watchdog...
2019-02-19 22:42:29.795955 INFO - client.go:62 StartClient: Start OVS Client
2019-02-19 22:42:29.885225 INFO - server.go:79 createService: Start generating initial gNMI config from OVS system source...
2019-02-19 22:42:29.951441 INFO - service.go:1182 StartService: Start gNXI Service
2019-02-19 22:42:29.952219 INFO - service.go:1188 StartService: Starting to listen
2019-02-19 22:42:29.952351 INFO - service.go:1194 StartService: Starting to serve gNXI
2019-02-19 22:42:38.271188 INFO - service.go:491 Capabilities: allowed a Capabilities request
2019-02-19 22:42:38.271540 INFO - service.go:504 Capabilities: Send Capability response to client: supported_models:<name:"openconfig-interfaces" organization:"OpenConfig working group" version:"2.0.0" > supported_models:<name:"openconfig-openflow" organization:"OpenConfig working group" version:"0.1.0" > supported_models:<name:"openconfig-platform" organization:"OpenConfig working group" version:"0.5.0" > supported_models:<name:"openconfig-system" organization:"OpenConfig working group" version:"0.2.0" > supported_encodings:JSON supported_encodings:JSON_IETF gNMI_version:"0.7.0"
2019-02-19 22:42:38.300135 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:38.300285 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616158300283711 update:<path:<elem:<name:"system" > elem:<name:"config" > elem:<name:"hostname" > > val:<string_val:"target.gnxi.lan" > > >
2019-02-19 22:42:38.329608 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:38.329878 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616158329876087 update:<path:<elem:<name:"components" > elem:<name:"component" key:<key:"name" value:"os" > > elem:<name:"state" > elem:<name:"description" > > val:<string_val:"2.9.2" > > >
2019-02-19 22:42:38.359250 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:38.359489 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616158359486853 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.2" > > >
2019-02-19 22:42:38.388493 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:38.389839 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616158389836628 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6653 > > >
2019-02-19 22:42:38.418856 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:38.419039 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616158419037234 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:21 > > >
2019-02-19 22:42:38.447727 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:38.448002 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616158448000256 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:16 > > >
2019-02-19 22:42:38.483561 INFO - service.go:980 Reboot: allowed a Reboot request
2019-02-19 22:42:38.483664 INFO - service.go:990 Reboot: Send Reboot response to client:
2019-02-19 22:42:41.664241 INFO - client.go:79 StopClient: Stop OVS Client
2019-02-19 22:42:41.664282 INFO - client.go:62 StartClient: Start OVS Client
2019-02-19 22:42:41.664970 INFO - service.go:1201 StopService: Stop gNXI Service
2019-02-19 22:42:41.695177 INFO - server.go:79 createService: Start generating initial gNMI config from OVS system source...
2019-02-19 22:42:41.707586 INFO - service.go:1182 StartService: Start gNXI Service
2019-02-19 22:42:41.707639 INFO - service.go:1188 StartService: Starting to listen
2019-02-19 22:42:41.707683 INFO - service.go:1194 StartService: Starting to serve gNXI
2019-02-19 22:42:48.519107 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:48.519288 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616168519285865 update:<path:<elem:<name:"system" > elem:<name:"config" > elem:<name:"hostname" > > val:<string_val:"target.gnxi.lan" > > >
2019-02-19 22:42:48.550731 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:48.550856 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616168550854665 update:<path:<elem:<name:"components" > elem:<name:"component" key:<key:"name" value:"os" > > elem:<name:"state" > elem:<name:"description" > > val:<string_val:"2.9.2" > > >
2019-02-19 22:42:48.579898 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:48.580093 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616168580091137 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.2" > > >
2019-02-19 22:42:48.614391 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:48.614567 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616168614565444 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6653 > > >
2019-02-19 22:42:48.651481 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:48.651715 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616168651713657 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:29 > > >
2019-02-19 22:42:48.688920 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:48.689119 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616168689117260 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:29 > > >
2019-02-19 22:42:48.729953 INFO - service.go:1140 GetCertificates: allowed a GetCertificates request
2019-02-19 22:42:48.729987 INFO - service.go:1146 GetCertificates: Send GetCertificates response to client: certificate_info:<certificate_id:"c5e5a1cb-8e1f-43c1-be4a-ab8e513fc667" certificate:<type:CT_X509 certificate:"-----BEGIN CERTIFICATE-----\nMIIFRjCCAy4CCQCkHpna8OHQHTANBgkqhkiG9w0BAQsFADBjMQswCQYDVQQGEwJB\nVDEPMA0GA1UECAwGVmllbm5hMQ8wDQYDVQQHDAZWaWVubmExDTALBgNVBAoMBFRl\nc3QxDTALBgNVBAsMBFRlc3QxFDASBgNVBAMMC2NhLmdueGkubGFuMB4XDTE5MDIx\nOTIxNDA1NFoXDTE5MDMyMTIxNDA1NFowZzELMAkGA1UEBhMCQVQxDzANBgNVBAgM\nBlZpZW5uYTEPMA0GA1UEBwwGVmllbm5hMQ0wCwYDVQQKDARUZXN0MQ0wCwYDVQQL\nDARUZXN0MRgwFgYDVQQDDA90YXJnZXQuZ254aS5sYW4wggIiMA0GCSqGSIb3DQEB\nAQUAA4ICDwAwggIKAoICAQDSiBNm5rG1d+FU07ZvC6pNfEvkLwetLKszEBpXZHdP\nftcksPaFE6kYhJz/Y48qsNaHuCWaxSAIpMJL3i1VMoVos4iuvYfiu1qiRMAllSLA\nQMnNg7XdxjVKstL9MD1mBX2ULHD0zgEKt9zjQHdop2Ou/9tGZaSjide0yoZSUCVm\nbq2+EEiiGc6xUYYE/SX58WDyasWrJLQwNBRkBtpb5UeXth1f+g1QN6gdTCwVJ5kY\nxktvrEjsprjxfTg7E5sWl1aC87CBltgNwknSF0JiysY9W1I0vA3MV5ifVEiaZYY4\n01sq1xjwMRxgPIHNMvznFkY/J8SiYB7OtilKW3iCDccmgAPlFkTSm52+ARRSQiz7\nJAyxAow2KRrDDFSsk6OqKKqhalXRzLSP9rDHofGFCqJ9U1cWJmXEt/ffr5Ruo46j\n9fUWAgi2yfE/9jhvURbvCOX2xgJJQSzc6Od5YLpavkVC3fqHTVuXg++XGoX7b8Qy\nc1kRs0rYOK8l3sr8HM8z43cfpReWfYUcskC6M1y9Vvi58SB8VG/xozNYyElmfsLi\nOu8z9AzoLg4wuZ067XejVgUrD/t6ESkd+zLkomKiKVI4m87ZVvMPoZC+CGbVNDeh\nRrMeiAAymGPh6byKIrOCBIDM1McyC+A65+Ebi1kP8SpqFxLUUuItOdo0QcQ1YIqi\nAwIDAQABMA0GCSqGSIb3DQEBCwUAA4ICAQBzg5XsOAvpb4QAfopV5skEQscLh+l4\nCmZMeB7D5eJfJQPXqkHvh9QQlIBQa+ejB4jFzMZL0xLYILDZI79x0b/GFlxQB1eC\n6MbLyP9Hc9a4NTuK7doHxw7jJ3Mskhi2PjvUKedreNcRMt5wskyfdqEoEZZCDT4v\npC2UgCDCZ46d1NmGWgvGzj6fFxvNnxrOhRQzGd1jGDUJktls4dg+uAn5RRuFJURH\niH+OE69QeVB3no0z0svLxXX2+qpYlsG5Ci712Q3CY3+KbETkSlVebrH8Ffo1UTmC\nhl+qWdwJN1saBS6lcSZ3ddY7TiP7Mjh8kQqxptEwPy45ZPj/a+UloIoL+SEH/yYd\naUZR8V7lLD+LIJhl3FX8WSadg3nxJ474XXFM4gwzhrrlmpY3f8A5oVE94tS5TwtA\nqsnWsBy4/QFrFauU9+EI0AGKX29OMwgLIHSbRF+yxl8ibj0WlE+tY8BxNCRN+aGr\nA3tcBG0DKbeBG3+Za2NvdnwTPib5HSdeQXSab9OcOI0MViW89jtxeu9nkC0Z7lzx\n7QHfBgUTEfJEhCNkNzGDs254rcbVtI2WFWonPz2iEolQdxkJVBkFsU1R7uj3yynF\nLOCqSlbYsPw457++iSmGg5FCLjY8T24YBasuCKPiNyHQD66C+UpvRfPpWQIozE8P\nKL2X7hqQqcdZIw==\n-----END CERTIFICATE-----\n" > modification_time:1550616149795568304 >
2019-02-19 22:42:48.761158 INFO - service.go:1032 Rotate: allowed a Rotate request: generate_csr:<csr_params:<type:CT_X509 min_key_size:4096 key_type:KT_RSA common_name:"target.gnxi.lan" country:"AT" state:"Vienna" city:"Vienna" organization:"Test" organizational_unit:"Test" > >
2019-02-19 22:42:51.537528 INFO - manager.go:192 ActivatePackage: Cert package d7f58600-4b8e-4260-be3d-ff1641e1c8e9 is now active
2019-02-19 22:42:54.702716 INFO - client.go:79 StopClient: Stop OVS Client
2019-02-19 22:42:54.702776 INFO - client.go:62 StartClient: Start OVS Client
2019-02-19 22:42:54.743864 INFO - service.go:1140 GetCertificates: allowed a GetCertificates request
2019-02-19 22:42:54.743899 INFO - service.go:1146 GetCertificates: Send GetCertificates response to client: certificate_info:<certificate_id:"d7f58600-4b8e-4260-be3d-ff1641e1c8e9" certificate:<type:CT_X509 certificate:"-----BEGIN CERTIFICATE-----\nMIIFlDCCA3ygAwIBAgIRAPoatp+OpfB+gIgS6QhImF0wDQYJKoZIhvcNAQELBQAw\nYzELMAkGA1UEBhMCQVQxDzANBgNVBAgMBlZpZW5uYTEPMA0GA1UEBwwGVmllbm5h\nMQ0wCwYDVQQKDARUZXN0MQ0wCwYDVQQLDARUZXN0MRQwEgYDVQQDDAtjYS5nbnhp\nLmxhbjAeFw0xOTAyMTkyMjQyNTFaFw0yMDAyMTkyMjQyNTFaMEUxCzAJBgNVBAYT\nAkFUMQ0wCwYDVQQKEwRUZXN0MQ0wCwYDVQQLEwRUZXN0MRgwFgYDVQQDEw90YXJn\nZXQuZ254aS5sYW4wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDPMSuX\njfRe+6egfPFUEI4nt0wL9H4qKPtvjfec5n50CpTT8O5OPo1DugMSf5CX8WiXrxPs\nkhbmch2Z+U8CtaaGr/5qbJZv0c5+1KbeS3wD326ieuwMmKh9NPeNuxlOWvPEiO/7\nr38fNHtZ7z0PlPHfobCMMORxK2kADOzDw1+6WCcB0aetluxVzk0vy4nTatgrRhCV\nl1ZAJgf23niIEe670rW3YJ2YUfLp51g1QaUsU6/3nvANZt0Kod7Kq0gTemJvxSJZ\noFVwvchsuhLrwiocZMdP/abtPj/isneSDN2999EHV9snMjrDM6P5bSYJaQbLKgbs\nqZIkF9uvrsvfQeD33opofZH5Vhj2Onojpuf9QkzPr5UA/Cac2cIkDK2ZofquOn0S\nmmu3NKJPhmJRyKxSu500koOUBlYrsqu4G6M2/hu40lZRG7OUxFHo5i5o+Hjbu05o\nYuxl//H1dRepF18auA0ysfr+xzjSlv+7j8zBiu/rzKCbEM39JhwK/Y/IVdtvrW5g\nNJvyTCE1DwjfySvjKpWR6WGh/eroZne441xBh2lUgHsrS0qhBYlREDRfZ+bpRD24\nZI/QZESSSxDxiSkzvDgLM/T6aFZedaT5cc3Npr13f6IFURQFeWwBARofV/kgUttF\nt+q3W7/oR6fsTdpNCnEFtlAc3Zyb1In2y5uZsQIDAQABo2EwXzAOBgNVHQ8BAf8E\nBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMBMA8GA1UdEwEB/wQF\nMAMCAQUwHQYDVR0OBBYEFCaII2uz3XY0xYLdlzJioQJ9AN6JMA0GCSqGSIb3DQEB\nCwUAA4ICAQAzaazKQ9qlNjCftEhCn5g295a15T3OiooP8Ab/gOVOYxi9mffNaJVF\n0ltuYfpPCHEs2jx3CDldgkhpjehWzAt2pCDpmlskIpqiZz7cUhnSLLaiL4T6KA4M\nOCAsMS18LMLIlG+vodDlK6y77tUdhtkL9RgRhUM0KV5lh3V96b4Lippyr1SIu2tC\n2nbakyCAQRYjl1XuS2oMSs3UawT81pzECDOlWfmisN/qWLJsQiY2dtkROejGuRAm\nrdIYQ53NZ02bCmX2P8GoUwD9SVRGeTK+Y6+t+7Yn9QXIdCADqo6Is0Gxa72dc45Y\n1o1tNwWhCErUuoTpvFBcxfIcpoaQkCE5i/QDKSX4si99kNWwZ/dLpvHHFUIrxRmG\nQSlOUQpSOn+iQxOhe9lo6YUADkBw7PfWS5u5Yddss6Wr9hgyI/V8ZVF+Iojt1UkW\nJ6vzyVS5VIA5awTfsN7EBrvor6ZIiaTuKeifzo6BkfbuvShTRLzUOaSW+8JSl5xs\nHfPVGHXwdHuj3alKklb+JG0G9QuSILvcJ+rp1gOuCEsVFO6xTsHEXt6LhpxwOguf\neHFqJew7HEwrxMZoN6oSEclQxiA8NcoVOpvTNyuL/o5NobSqHasp1k8ll0MaCVA3\n1i+L6T9JPF2h0mzLEBwK3ld+XjY0e5QF0DjNkpHuAX26mqnXz2F83Q==\n-----END CERTIFICATE-----\n" > modification_time:1550616171533236982 >
2019-02-19 22:42:54.794847 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:54.795034 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616174795032701 update:<path:<elem:<name:"system" > elem:<name:"config" > elem:<name:"hostname" > > val:<string_val:"target.gnxi.lan" > > >
2019-02-19 22:42:54.826143 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:54.826578 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616174826574966 update:<path:<elem:<name:"components" > elem:<name:"component" key:<key:"name" value:"os" > > elem:<name:"state" > elem:<name:"description" > > val:<string_val:"2.9.2" > > >
2019-02-19 22:42:54.872726 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:54.872902 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616174872900572 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.2" > > >
2019-02-19 22:42:54.904851 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:54.905266 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616174905261172 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6653 > > >
2019-02-19 22:42:54.935130 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:54.935346 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616174935343750 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:38 > > >
2019-02-19 22:42:54.964279 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:54.964681 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616174964679288 update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:34 > > >
2019-02-19 22:42:55.023418 INFO - service.go:639 Set: allowed a Set request
2019-02-19 22:42:55.028726 INFO - client.go:212 SyncChangesToRemote: target is in inconsistent state with OVS device, syncing Controller
2019-02-19 22:42:55.036174 INFO - service.go:694 Set: Send Set response to client: response:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > op:UPDATE >
2019-02-19 22:42:55.074898 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:55.075160 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616175075158280 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.3" > > >
2019-02-19 22:42:55.106035 INFO - service.go:639 Set: allowed a Set request
2019-02-19 22:42:55.111517 INFO - client.go:212 SyncChangesToRemote: target is in inconsistent state with OVS device, syncing Controller
2019-02-19 22:42:55.139021 INFO - service.go:694 Set: Send Set response to client: response:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > op:UPDATE >
2019-02-19 22:42:55.170313 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:55.170486 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616175170484527 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"address" > > val:<string_val:"172.18.0.2" > > >
2019-02-19 22:42:55.201126 INFO - service.go:639 Set: allowed a Set request
2019-02-19 22:42:55.208394 INFO - client.go:212 SyncChangesToRemote: target is in inconsistent state with OVS device, syncing Controller
2019-02-19 22:42:55.216583 INFO - service.go:694 Set: Send Set response to client: response:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > op:UPDATE >
2019-02-19 22:42:55.250020 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:55.250199 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616175250198034 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6654 > > >
2019-02-19 22:42:55.280240 INFO - service.go:639 Set: allowed a Set request
2019-02-19 22:42:55.286661 INFO - client.go:212 SyncChangesToRemote: target is in inconsistent state with OVS device, syncing Controller
2019-02-19 22:42:55.305546 INFO - service.go:694 Set: Send Set response to client: response:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > op:UPDATE >
2019-02-19 22:42:55.344623 INFO - service.go:516 Get: allowed a Get request
2019-02-19 22:42:55.344792 INFO - service.go:627 Get: Send Get response to client: notification:<timestamp:1550616175344789865 update:<path:<elem:<name:"system" > elem:<name:"openflow" > elem:<name:"controllers" > elem:<name:"controller" key:<key:"name" value:"main" > > elem:<name:"connections" > elem:<name:"connection" key:<key:"aux-id" value:"0" > > elem:<name:"config" > elem:<name:"port" > > val:<uint_val:6653 > > >
2019-02-19 22:42:55.374919 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > > mode:ONCE encoding:JSON_IETF >
2019-02-19 22:42:55.375088 INFO - service.go:877 subscribeOnce: serving subscribe ONCE
2019-02-19 22:42:55.375286 INFO - service.go:885 subscribeOnce: Send Subscribe ONCE response to client: update:<timestamp:1550616175375283398 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:38 > > >
2019-02-19 22:42:55.403714 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > > mode:ONCE encoding:JSON_IETF >
2019-02-19 22:42:55.403923 INFO - service.go:877 subscribeOnce: serving subscribe ONCE
2019-02-19 22:42:55.404212 INFO - service.go:885 subscribeOnce: Send Subscribe ONCE response to client: update:<timestamp:1550616175404202730 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:34 > > >
2019-02-19 22:42:55.433328 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > > mode:POLL encoding:JSON_IETF >
2019-02-19 22:42:55.433797 INFO - service.go:897 subscribePoll: serving subscribe POLL
2019-02-19 22:42:55.434073 INFO - service.go:903 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1550616175434067630 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:38 > > >
2019-02-19 22:43:00.436913 INFO - service.go:940 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1550616180436909598 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:43 > > >
2019-02-19 22:43:05.438593 INFO - service.go:940 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1550616185438580667 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:49 > > >
2019-02-19 22:43:05.467723 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > > mode:POLL encoding:JSON_IETF >
2019-02-19 22:43:05.467795 INFO - service.go:897 subscribePoll: serving subscribe POLL
2019-02-19 22:43:05.467991 INFO - service.go:903 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1550616185467988694 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:45 > > >
2019-02-19 22:43:10.469971 INFO - service.go:940 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1550616190469966746 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:51 > > >
2019-02-19 22:43:15.471179 INFO - service.go:940 subscribePoll: Send Subscribe POLL response to client: update:<timestamp:1550616195470905116 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:57 > > >
2019-02-19 22:43:15.508762 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > > encoding:JSON_IETF >
2019-02-19 22:43:15.508832 INFO - service.go:952 subscribeStream: serving subscribe STREAM
2019-02-19 22:43:16.712154 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1550616196712151135 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:66 > > >
2019-02-19 22:43:21.714061 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1550616201714058198 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:71 > > >
2019-02-19 22:43:26.718358 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1550616206718354950 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:75 > > >
2019-02-19 22:43:26.748723 INFO - service.go:721 Subscribe: allowed Subscribe request: subscribe:<prefix:<target:"target.gnxi.lan" > subscription:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > > encoding:JSON_IETF >
2019-02-19 22:43:26.749057 INFO - service.go:952 subscribeStream: serving subscribe STREAM
2019-02-19 22:43:31.719970 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1550616211719966660 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"in-pkts" > > val:<uint_val:81 > > >
2019-02-19 22:43:36.719072 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1550616216718959318 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:84 > > >
2019-02-19 22:43:41.724193 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1550616221724190321 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:89 > > >
2019-02-19 22:43:46.724480 INFO - service.go:963 subscribeStream: Send Subscribe STREAM response to client: update:<timestamp:1550616226724477399 prefix:<target:"target.gnxi.lan" > update:<path:<elem:<name:"interfaces" > elem:<name:"interface" key:<key:"name" value:"sw1-eth1" > > elem:<name:"state" > elem:<name:"counters" > elem:<name:"out-pkts" > > val:<uint_val:94 > > >
```

## Generate Go Bindings for Open vSwitch
go run vendor/github.com/openconfig/ygot/generator/generator.go -path=yang -output_file=oc.go -package_name=openvswitch -generate_fakeroot -fakeroot_name=device -compress_paths=true -exclude_modules=ietf-interfaces yang/openconfig-interfaces.yang