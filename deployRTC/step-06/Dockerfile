FROM node:10-alpine

WORKDIR /deployRTC/step06

COPY package.json /deployRTC/step06

RUN npm install

COPY . /deployRTC/step06

CMD node index.js

EXPOSE 8080
