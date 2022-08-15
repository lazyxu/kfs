if (process.env.REACT_APP_PLATFORM === 'web') {
  module.exports = require('./config.web.js');
} else {
  module.exports = require('./config.electron.js');
}
