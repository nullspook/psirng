#!/bin/bash
set -e

apt update
apt -y install g++ cmake libusb-1.0-0-dev libconfuse-dev

ln -s $(find /usr/lib -name libusb-1.0.a -print -quit) /usr/local/lib/libusb-1.0.a
ln -s $(find /usr/lib -name libusb-1.0.so -print -quit) /usr/local/lib/libusb-1.0.so

cd /app/3rdparty/libqwqngx/3rdparty/libftdi1
mkdir -p build && cd build
cmake ..
make -j$(nproc)
make install
ln -sf /usr/local/include/libftdi1/ftdi.h /usr/local/include/ftdi.h

cd /app/3rdparty/libqwqngx/3rdparty/libqwqng
mkdir -p build && cd build
cmake ..
make -j$(nproc)
make install

cd /app/3rdparty/libqwqngx
mkdir -p build && cd build
cmake ..
make -j$(nproc)
make install
ldconfig
