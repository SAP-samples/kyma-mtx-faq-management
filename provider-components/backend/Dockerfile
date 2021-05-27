FROM node:12.18.1 as builder
USER root
RUN apt-get update && apt-get install -y python3
RUN npm i -g @sap/cds-dk
# create work directory and copy sources  
WORKDIR /cap
COPY . .
# install dependencies
RUN npm install
#execute build
RUN cds build --production
FROM node:12.18.1
# Set production environment
ENV NODE_ENV=production
# Create app directory
WORKDIR /usr/src/app
# Copy from builder and install Install app dependencies
COPY --from=builder /cap/gen/srv .
RUN npm install
EXPOSE 4004
CMD [ "npm", "start" ]