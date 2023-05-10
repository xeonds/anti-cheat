FROM ubuntu:latest

ENV OPENAI_API_KEY ''

RUN apt-get update && apt-get install -y --no-install-recommends nodejs \
    && apt-get purge --auto-remove \
    && rm -rf /tmp/* /var/lib/apt/lists/*

WORKDIR /app
COPY package.json /app/
COPY index.js /app/
COPY config.js /app/
RUN npm i --registry=https://registry.npm.taobao.org/ \
    && rm -fr /tmp/* ~/.npm

ENTRYPOINT [ "node", "index.js" ]
