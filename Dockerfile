FROM alpine

ENV GOPATH=$HOME/go
ENV APP_HOME=$GOPATH/src/github.com/softleader/deployer

WORKDIR $APP_HOME

#COPY main.go $APP_HOME
#COPY cmd/ $APP_HOME/cmd/
#COPY datamodels/ $APP_HOME/datamodels/
#COPY services/ $APP_HOME/services/
#COPY web/ $APP_HOME/web/

COPY main $APP_HOME

RUN apk update \
    && apk --no-cache add git nodejs-npm go \
    && rm -rf /var/cache/apk/* \
    && git config --global user.name "r&d" \
    && git config --global user.email rd@softleader.com.tw \
    && npm install softleader/container-yaml-generator -g \
    && npm install softleader/git-package-manager -g

EXPOSE 8080

CMD ["sh", "-c", "./main"]