FROM node:18-alpine as build

WORKDIR /app

COPY . .
RUN npm install

EXPOSE 8080
CMD ["npm", "run", "start"]
