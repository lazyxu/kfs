const { override, useBabelRc, babelInclude } = require('customize-cra');

const path = require('path');

module.exports = override(
  useBabelRc(),
  babelInclude([path.normalize(`${path.resolve()}/src`)]),
);
