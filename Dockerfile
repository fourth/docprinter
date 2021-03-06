FROM ubuntu:14.04
RUN apt-get update && apt-get install -y git \
    mercurial \
    curl \
    ca-certificates \
    build-essential \
    libfreetype6 \
    libfontconfig1 \
    bzr \
    --no-install-recommends

RUN	curl -sSL https://golang.org/dl/go1.3.1.src.tar.gz | tar -v -C /usr/local -xz
ENV	PATH	/usr/local/go/bin:$PATH
ENV	GOPATH	/go:/go/src/github.com/fourth/docprinter/Godeps/_workspace
ENV     PATH /go/bin:$PATH
RUN	cd /usr/local/go/src && ./make.bash --no-clean 2>&1

RUN     curl -sSL https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-1.9.7-linux-x86_64.tar.bz2 | tar -jxv
RUN     cp phantomjs-1.9.7-linux-x86_64/bin/phantomjs /usr/local/bin

RUN     apt-get install -y fontconfig

COPY	. /go/src/github.com/fourth/docprinter
WORKDIR /go/src/github.com/fourth/docprinter
RUN     cp -R fonts/* /usr/share/fonts/truetype/
RUN     fc-cache -f -v

RUN     go build

ENTRYPOINT ["./docprinter"]
