FROM golang:1.16rc1-buster

# RUN apt-get update && apt-get -y upgrade && apt-get install -y libcairo2-dev gcc pkg-config libjpeg-turbo-dev libpng-dev libz-dev cmake libboost-dev libglib-dev

# RUN cd /popper-build/build && 
# cmake ../ \
#   -DCMAKE_BUILD_TYPE=Release \
#   -DCMAKE_INSTALL_PREFIX:PATH=/poppler \
#   -DCMAKE_INSTALL_LIBDIR=/usr/lib \
#   -DENABLE_UNSTABLE_API_ABI_HEADERS=OFF \
#   -DBUILD_SHARED_LIBS=OFF

RUN apt-get update && apt-get install gcc g++

RUN wget -O- https://github.com/ashutoshvarma/libxpdf/releases/download/v0.1.3/libxpdf-4.02.linux-gcc.x64.zip | unzip -
