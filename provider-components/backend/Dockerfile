FROM node:14-alpine as builder
USER root
RUN npm i -g @sap/cds-dk
# create work directory and copy sources  
WORKDIR /cap
COPY . .
# install dependencies
RUN npm install
#execute build
RUN cds build --production

FROM node:14-alpine
# Set production environment
ENV NODE_ENV=production
# Create app directory
WORKDIR /usr/src/app
# Copy from builder and install Install app dependencies
COPY --from=builder /cap/gen/srv .
RUN npm install
EXPOSE 4004
CMD [ "npm", "start" ]