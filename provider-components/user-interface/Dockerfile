FROM node:10

# Create app directory
WORKDIR /app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY approuter/package.json ./

RUN npm install
# If you are building your code for production
# RUN npm ci --only=production

# Bundle app source
COPY approuter .

# Expose over port 8017
EXPOSE 5000

# Execute the node command to build and run your application.
CMD [ "node", "approuter-start.js" ]