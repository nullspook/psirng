psirng
======

An HTTP and gRPC API server for interacting with Quantum World Corporation
(QWC) / ComScire QRNGs (Quantum Random Number Generators).

This server eliminates driver-level buffering by flushing the internal buffer
on every request, ensuring that the random numbers are freshly generated.

Running
-------

```
docker build -t psirng .
docker run -d --name psirng \
           --device /dev/bus/usb `# or --privileged for more stability` \
           -p 8443:8443 \
           -p 8080:8080 \
           -p 50051:50051 \
           -v /path/to/your/cert.pem:/etc/psirng/cert.pem `# highly recommended` \
           -v /path/to/your/key.pem:/etc/psirng/key.pem `# highly recommended` \
           psirng
```

Usage
-----

### HTTP API curl examples

```
$ curl 'http://localhost:8080/api/randbooleans?length=3'
[false,true,true]

$ curl 'http://localhost:8080/api/randbytes?length=32'
"G4xZkM87+0hyZP47gHhDm1DE+L11Qv2ZRK9ZjKKoVvI="

$ curl 'http://localhost:8080/api/randintegers?min=-5&max=10&length=10'
[9,6,9,9,7,-5,9,1,-1,6]

$ curl 'http://localhost:8080/api/randuniform?min=-5.0&max=10.0&length=3'
[-3.0741757164702896,9.293731260430167,6.09011301791039]

$ curl 'http://localhost:8080/api/randnormal?mean=0.0&stddev=1.0&length=3'
[0.4821670713280625,-1.5485153368469722,1.2983731290714078]

$ curl 'http://localhost:8080/healthz'
{"status":"OK","time":"2025-06-01T10:23:26.480337861Z"}
```

### gRPC API

Use [rng.proto](grpc/rng.proto) and gRPC's [health.proto](https://github.com/grpc/grpc-proto/blob/6565a1ba38af695ace7c3ce6e6ff837ee87d4c10/grpc/health/v1/health.proto)
to generate the client code for your preferred language. Or use [libpsirngclient](https://github.com/nullspook/libpsirngclient)
if you are using C or C++.

License
-------

    Copyright (C) 2025 NullSpook

    psirng is free software: you can redistribute it and/or modify it under the
    terms of the GNU Affero General Public License as published by the Free
    Software Foundation, either version 3 of the License, or (at your option)
    any later version.

    psirng is distributed in the hope that it will be useful, but WITHOUT ANY
    WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
    FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License for
    more details.

    You should have received a copy of the GNU Affero General Public License
    along with psirng.  If not, see <https://www.gnu.org/licenses/>.