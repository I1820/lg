# IoT Load Generator :fire:

[![Drone (cloud)](https://img.shields.io/drone/build/toskatok/lg.svg?style=flat-square)](https://cloud.drone.io/toskatok/lg)
[![GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

[![](https://images.microbadger.com/badges/image/toskatok/lg.svg)](https://microbadger.com/images/toskatok/lg "Get your own image badge on microbadger.com")

**Sometimes we tell lies, sometimes we prove we don't lie. Let's prove ourselves.**

## Introduction
This application gives a way for creating a load with MQTT, HTTP ...
MQTT is a messaging protocol and various platforms work with it so
creating a load with MQTT creates a way for testing different platforms.

LG has an executable version that you can run on your system with ease.
Besides, it has a server version which has APIs that are listed [here](https://app.swaggerhub.com/apis/I1820/i1820-lg/1.0.0).
In the server version reports about the status of each test's instance are emitted on
the socket.io room which has the name I1820,
with an event that has the name of its instance.

## Transport
LG have an awesome way for customizing portocol named `Transport`.
Transports give you a way for sending your generated data over *your protocol* with ease.

## Generator
LG is fully customizable so you can generate data with
your own structure and publish it on your own topic.
For this LG has the `Generator` interface that is defined in `generator/`.

## Running
You can run LG execuable version with following syntax:

```sh
lg --destination http://127.0.0.1:1883 --rate 1ms
```

In `config.yml` you can specifies generator configuration.
If your destination has scheme `http://` your transport will be HTTP
and if it has `mqtt://` your transport will be MQTT.

## Use Cases :male_detective:
### Set [loraserver.io](https://www.loraserver.io/) on fire
### ABP
With this load generator, we try to put a load on loraserver.io
and we get the following results.

```yml
generator:
  name: lora
  info:
    gateway:
      mac: "b827ebffff70c80a"
    keys:
      networkSkey: "DB56B6C3002A4763A79E64573C629D97"
      applicationSKey: "94B49CD7BC621BC46571D019640804AA"
    device:
      addr: "26011CF6"
messages:
  - 100: 6750
    101: 6606
    lat: 10
    lng: 10
```

| Interval | Status        |
|:---------|:-------------:|
| 1s       | Success       |
| 100ms    | Success       |
| 10ms     | Success       |
| 1ms      | Fail          |


### Set [I1820](https://i1820.org) on fire
#### TTN (over HTTP)
With this load generator, we try to put a load on I1820 TTN Integration module in link component
and we get the following results.

```yml
generator:
  name: ttn
  info:
    applicationName: fan
    applicationID: 5ba3f19c87a142b0a840fae0
    devEUI: 000AE31955C049FC
    deviceName: agrinode
token: ttnIStheBEST
messages:
  - count: "{{.Count}}"
```

These results show generated parsed information ratio (number of parsed data / number of received data)
with send interval of data:

| Interval | Accept Ratio  |
|:---------|:-------------:|
| 1s       | 100%          |
| 100ms    | 100%          |
| 10ms     | 100%          |
| 1ms      | 100%          |
| 100us    | 100%          |
| 10us     | 100%          |

Please consider that HTTP requests cannot get their response in `100us` or lower interval so in these intervals there is no
difference with `1ms` interval.

#### Fanco (over MQTT)
With this load generator, we try to put a load on I1820 MQTT service in link component
and we get the following results.

```yml
generator:
  name: fanco
  info:
    thingID: 5bbd104cefe940cb57dfeb76
token: 1BLzO2YYB1jH91pRB0cpeIdPMsM
messages:
  - count: "{{.Count}}"
```

These results show generated parsed information ratio (number of parsed data / number of received data)
with send interval of data:

| Interval | Accept Ratio  |
|:---------|:-------------:|
| 1s       | 100%          |
| 100ms    | 100%          |
| 10ms     | 100%          |
| 1ms      | 100%          |
| 100us    | 52.3809524%   |
| 10us     | 47.6190476%   |

#### Test System Specification

- 4 Core Intel(R) Core(TM) i7-5930K CPU @ 3.50GHz
- 8 Gb of RAM
- SSD
