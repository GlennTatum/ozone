FROM node:trixie-slim

WORKDIR /app

COPY . /app/

RUN npm install -g @angular/cli && npm install

EXPOSE 4200

ENTRYPOINT [ "/usr/local/bin/ng", "serve", "--host", "0.0.0.0" ]