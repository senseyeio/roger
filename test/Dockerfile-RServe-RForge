FROM r-base:3.5.1

RUN apt-get update

RUN apt-get install -y --no-install-recommends \
              libxml2-dev \
              libcurl4-gnutls-dev \
              libssl-dev

RUN R -e 'install.packages("Rserve",,"http://rforge.net/")'

RUN echo "port 6311" >> /etc/Rserv.conf
RUN echo "remote enable" >> /etc/Rserv.conf

RUN echo "port 6312" >> /etc/Rserv-secure.conf
RUN echo "remote enable" >> /etc/Rserv-secure.conf
RUN echo "auth required" >> /etc/Rserv-secure.conf
RUN echo "pwdfile /etc/Rserve.pwd" >> /etc/Rserv-secure.conf
RUN echo "roger testpassword" >> /etc/Rserve.pwd

COPY . /usr/local/src/senseyeio
WORKDIR /usr/local/src/senseyeio

CMD nohup R < test.r --no-save & nohup R < test-secure.r --no-save

EXPOSE 6311 6312
