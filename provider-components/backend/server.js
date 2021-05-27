const cds = require('@sap/cds');

cds.on('bootstrap', (app) => {
  cds.mtx.in(app); // serve cds-mtx APIs
});

// Delegate bootstrapping to built-in server.js
module.exports = cds.server
