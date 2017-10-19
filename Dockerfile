FROM gcc:4.6

RUN apt-get update && apt-get install -y curl ed make bison git mercurial

RUN curl -Ls https://github.com/golang/go/archive/weekly.2010-12-22.tar.gz | tar -xzvf - && mv go-weekly.2010-12-22 /usr/local/go

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV GOBIN  $GOPATH/bin

RUN mkdir -p $GOPATH $GOBIN && cd $GOROOT/src && ./make.bash

ENV PATH $GOBIN:$PATH

COPY    . $GOPATH/src/github.com/creack/goshell
WORKDIR $GOPATH/src/github.com/creack/goshell

RUN make

CMD ./gosh
