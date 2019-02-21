FROM softleader/docker
MAINTAINER softleader.com.tw

RUN apk update \
    && apk --no-cache add git nodejs openssh \
    && rm -rf /var/cache/apk/* \
    && git config --global user.name "SoftLeader" \
    && git config --global user.email "support@softleader.com.tw"

ENV APP_HOME=/deployer
ENV WORKSPACE=/workspace
ENV CMD_GPM=$APP_HOME/node_modules/git-package-manager/index.js
ENV CMD_GEN_YAML=$APP_HOME/node_modules/container-yaml-generator/index.js
ENV PORT=80

WORKDIR $APP_HOME

COPY docker-compose.yml /
COPY node_modules/ $APP_HOME/node_modules/
COPY static/ $APP_HOME/static/
COPY templates/ $APP_HOME/templates/
COPY build/main $APP_HOME/

EXPOSE $PORT

CMD ["sh", "-c", "/deployer/main -workspace=$WORKSPACE -cmd.gpm=$CMD_GPM -cmd.gen-yaml=$CMD_GEN_YAML -port=$PORT"]
