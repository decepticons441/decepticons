FROM node:10-alpine

WORKDIR /deployRTC/convertToAWS

COPY ./package.json /deployRTC/convertToAWS

COPY . /deployRTC/convertToAWS

RUN npm install

EXPOSE 80

ENTRYPOINT ["node", "index.js"]


