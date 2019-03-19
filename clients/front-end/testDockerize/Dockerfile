FROM node:10-alpine

WORKDIR /deployRTC/testDockerize

COPY ./package.json /deployRTC/testDockerize

RUN npm install

COPY . /deployRTC/testDockerize

ENTRYPOINT ["node", "index.js"]

EXPOSE 80
