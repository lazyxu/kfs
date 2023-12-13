if (!import.meta.env.VITE_APP_PLATFORM || import.meta.env.VITE_APP_PLATFORM === 'web') {
  module.exports = require('./config.web.js');
} else {
  module.exports = require('./config.electron.js');
}
