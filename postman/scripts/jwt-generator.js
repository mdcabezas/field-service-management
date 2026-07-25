// jwt-generator.js — Pre-request script for Postman
// Generates JWT HS256 using crypto-js from Postman Sandbox
// Usage: Set in Pre-request Script of collection or folder

(function () {
  // Base64URL encode
  function base64url(str) {
    return btoa(str)
      .replace(/=/g, '')
      .replace(/\+/g, '-')
      .replace(/\//g, '_');
  }

  // HMAC-SHA256 using Postman's crypto-js
  function hmacSha256(message, secret) {
    var hash = CryptoJS.HmacSHA256(message, secret);
    return hash.toString(CryptoJS.enc.Base64)
      .replace(/=/g, '')
      .replace(/\+/g, '-')
      .replace(/\//g, '_');
  }

  // Check if token is still valid (not expired)
  var existingToken = pm.environment.get('jwt_token');
  if (existingToken) {
    try {
      var parts = existingToken.split('.');
      var payload = JSON.parse(atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')));
      var now = Math.floor(Date.now() / 1000);
      if (payload.exp && payload.exp > now + 60) {
        // Token still valid (with 60s buffer)
        return;
      }
    } catch (e) {
      // Invalid token, regenerate
    }
  }

  var secret = pm.environment.get('PGRST_JWT_SECRET');
  if (!secret) {
    console.error('PGRST_JWT_SECRET not set in environment');
    return;
  }

  var userId = pm.environment.get('employee_number') || '1001';
  var userRole = pm.environment.get('user_role') || 'admin';
  var expiresInMinutes = parseInt(pm.environment.get('jwt_ttl_minutes')) || 60;

  var header = { alg: 'HS256', typ: 'JWT' };
  var now = Math.floor(Date.now() / 1000);
  var claims = {
    sub: userId,
    role: userRole,
    iat: now,
    exp: now + (expiresInMinutes * 60)
  };

  var encodedHeader = base64url(JSON.stringify(header));
  var encodedClaims = base64url(JSON.stringify(claims));
  var signature = hmacSha256(encodedHeader + '.' + encodedClaims, secret);

  var jwt = encodedHeader + '.' + encodedClaims + '.' + signature;
  pm.environment.set('jwt_token', jwt);
  console.log('JWT generated for user: ' + userId + ' (role: ' + userRole + ', exp: ' + expiresInMinutes + 'min)');
})();
