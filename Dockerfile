FROM golang:1.19-bullseye

WORKDIR /app
COPY . .

# Install dependencies

RUN apt update \
    && apt -y install git bzip2 g++ cmake libusb-1.0-0-dev libconfuse-dev \
    && ln -s $(find /usr/lib -name libusb-1.0.a -print -quit) /usr/local/lib/libusb-1.0.a \
    && ln -s $(find /usr/lib -name libusb-1.0.so -print -quit) /usr/local/lib/libusb-1.0.so

RUN cd 3rdparty/libqwqngx/3rdparty/libftdi1/ \
    && mkdir build \
    && cd build \
    && cmake .. \
    && make \
    && make install \
    && ln -s /usr/local/include/libftdi1/ftdi.h /usr/local/include/ftdi.h

RUN cd 3rdparty/libqwqngx/3rdparty/libqwqng/ \
    && mkdir build \
    && cd build \
    && cmake .. \
    && make \
    && make install

RUN cd 3rdparty/libqwqngx/ \
    && mkdir build \
    && cd build \
    && cmake .. \
    && make \
    && make install \
    && ldconfig

# Build psirng

RUN go build -o /usr/local/bin/psirng ./cmd/psirng

# Use testing certificate if not provided

RUN mkdir -p /etc/psirng \
    && [ ! -f /etc/psirng/cert.pem ] && [ ! -f /etc/psirng/key.pem ] && cp /app/cert.pem /etc/psirng/cert.pem && cp /app/key.pem /etc/psirng/key.pem

# Cleanup

RUN apt -y purge --auto-remove git bzip2 g++ cmake libconfuse-dev \
    && apt clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

EXPOSE 8443
EXPOSE 8080
EXPOSE 50051

CMD ["psirng", "-cert=/etc/psirng/cert.pem", "-key=/etc/psirng/key.pem"]