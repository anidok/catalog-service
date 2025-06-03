const jwt = require('jsonwebtoken');

const secret = 'some-key';  // Same secret as in jwt-helper.sh

function generateToken() {
    const payload = {
        iss: 'kong-jwt-auth',
        exp: Math.floor(Date.now() / 1000) + 3600,
        sub: 'user125'
    };
    const token = jwt.sign(payload, secret);
    console.log('Generated JWT Token:');
    console.log(token);
}

function verifyToken(token) {
    try {
        const decoded = jwt.verify(token, secret);
        console.log('Token is valid. Decoded payload:');
        console.log(JSON.stringify(decoded, null, 2));
    } catch (err) {
        console.error('Token verification failed:', err.message);
    }
}

// Handle command line arguments
const command = process.argv[2];
const token = process.argv[3];

if (command === 'generate') {
    generateToken();
} else if (command === 'verify' && token) {
    verifyToken(token);
} else {
    console.log('Usage:');
    console.log('  node generate-jwt.js generate');
    console.log('  node generate-jwt.js verify <token>');
    process.exit(1);
} 