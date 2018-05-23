FROM openjdk:8

ENV SYMMETRIC_VERSION_MAJOR 3.9
ENV SYMMETRIC_VERSION_MINOR 3
ENV SYMMETRIC_VERSION ${SYMMETRIC_VERSION_MAJOR}.${SYMMETRIC_VERSION_MINOR}

RUN wget https://sourceforge.net/projects/symmetricds/files/symmetricds/symmetricds-${SYMMETRIC_VERSION_MAJOR}/symmetric-server-${SYMMETRIC_VERSION}.zip/download -O download.zip && \
    unzip download.zip && rm download.zip && mv symmetric-server-${SYMMETRIC_VERSION} /opt/symmetric && \
    mv /opt/symmetric/conf/log4j.xml /opt/symmetric/conf/log4j.xml.bak

ADD log4j.xml /opt/symmetric/conf/

CMD ["/opt/symmetric/bin/sym"]
WORKDIR /opt/symmetric/bin
