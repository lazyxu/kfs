const path = require("path");
const fs = require("fs");
const ModuleScopePlugin = require('react-dev-utils/ModuleScopePlugin');
const rewireBabelLoader = require("craco-babel-loader");

// helpers

const appDirectory = fs.realpathSync(process.cwd());
const resolveApp = relativePath => path.resolve(appDirectory, relativePath);

module.exports = {
  webpack: {
      configure: {
          target: 'electron-renderer'
      }
  },
  plugins: [
    {
      plugin: {
        overrideWebpackConfig: ({ webpackConfig, cracoConfig, pluginOptions, context: { env, paths } }) => {
          webpackConfig.resolve.plugins = webpackConfig.resolve.plugins.filter(plugin => !(plugin instanceof ModuleScopePlugin));
          return webpackConfig;
        },
      },
    },
    //This is a craco plugin: https://github.com/sharegate/craco/blob/master/packages/craco/README.md#configuration-overview
    {
      plugin: rewireBabelLoader,
      options: {
        includes: [resolveApp("../common")], //put things you want to include in array here
        // excludes: [/(node_modules|bower_components)/] //things you want to exclude here
        //you can omit include or exclude if you only want to use one option
      }
    },
  ]
}
