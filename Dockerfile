# Fritz!Box Spectrum Logger (https://github.com/c-mueller/fritzbox-spectrum-logger).
# Copyright (c) 2018 Christian Müller <cmueller.dev@gmail.com>.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 3.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
# General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.
FROM ubuntu:18.04

WORKDIR /fsl

VOLUME /fsl/data

ENV DB_PATH="/fsl/data/spectra.db"
ENV AUTOLAUNCH="false"
ENV ENDPOINT_URL=":80"

ENV FRITZ_ENDPOINT="192.168.178.1"
ENV FRITZ_USERNAME=""
ENV FRITZ_PASSWORD=""
ENV UPDATE_INTERVAL="60"
ENV SESSION_REFRESH_INTERVAL="3600"
ENV SESSION_REFRESH_ATTEMPTS="5"
ENV MAX_DOWNLOAD_FAILS="5"

ADD fritzbox-spectrum-logger /usr/bin/fsl

EXPOSE 80
CMD fsl server env -d