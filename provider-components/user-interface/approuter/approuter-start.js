const approuter = require('@sap/approuter');
const jwt_decode = require('jwt-decode');

const ar = approuter();

ar.beforeRequestHandler.use('/me', function (req, res) {
  if (!req.user) {
    res.statusCode = 403;
    res.end(`Missing JWT Token`);
  } else {
    res.statusCode = 200;
    var token = jwt_decode(req.user.token.accessToken);
    var assignedScopes = token.scope;

    res.end(JSON.stringify({
      id: req.user.id,
      name: req.user.name,
      assignedScopes: assignedScopes
    }));
  }
});
ar.start();
