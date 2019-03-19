FROM node:latest

WORKDIR /deployRTC/andAnotherVidChat

COPY ./package.json /deployRTC/andAnotherVidChat

RUN npm install

COPY . /deployRTC/andAnotherVidChat

EXPOSE 80

ENTRYPOINT ["node", "server.js"]


