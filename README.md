psirng
======

An HTTP and GRPC API server for interacting with Quantum World Corporation
(QWC) / ComScire / MindEnabled QRNGs for psi experiments.

Running
-------

```
docker build -t psirng .
docker run -d --name psirng \
           --privileged \
           -p 8443:8443 \
           -p 8080:8080 \
           -p 50051:50051 \
           -v /path/to/your/cert.pem:/etc/psirng/cert.pem `# highly recommended` \
           -v /path/to/your/key.pem:/etc/psirng/key.pem `# highly recommended` \
           psirng
```

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