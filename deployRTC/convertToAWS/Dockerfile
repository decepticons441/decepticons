FROM node:10-alpine

WORKDIR /deployRTC/convertToAWS

COPY package.json /deployRTC/convertToAWS

RUN npm install

COPY . /deployRTC/convertToAWS

EXPOSE 80

CMD node index.js

