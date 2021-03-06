# IoT Sample

A simple service in Go to be used in demos.

[this blog post](http://modocache.svbtle.com/restful-go) was taken as a base

## Dependencies

It has a Kafka and MongoDB dependency in runtime.

## Running the Server

Make sure you have Kafka and MongoDB installed and running somewhere.

```
src/iot-demo/ $ go install
src/iot-demo/ $ iot-demo
[martini] listening on :3000 (development)
```

## Querying Data

### Register a new device

```
curl -i -X POST \
    -H "Content-Type: application/json" \
    -d '{"UniqueDeviceId": "ABCDEFG", "DeviceType": "10"}' \
    172.18.44.141:3000/api/device/RegisterDevice
```

### Get All Devices

```
$ curl -i localhost:3000/api/device/list
```
