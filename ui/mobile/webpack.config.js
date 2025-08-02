const createExpoWebpackConfigAsync = require('@expo/webpack-config');

module.exports = async function (env, argv) {
  const config = await createExpoWebpackConfigAsync(env, argv);

  // 1. 解决图标模块别名问题
  config.resolve.alias = {
    ...config.resolve.alias,
    '@react-native-vector-icons/material-design-icons': '@expo/vector-icons/MaterialCommunityIcons',
  };

  config.ignoreWarnings = [
    {
      module: /node_modules\/@microsoft\/fetch-event-source/
    },
    {
      module: /node_modules\/react-native-fetch-event-source/
    },
    { message: /Failed to parse source map/ } // 匹配所有类似警告
  ];

  // 2. 忽略 source-map-loader 对特定模块的警告
  config.module.rules = config.module.rules.map(rule => {
    if (rule.use?.loader === 'source-map-loader') {
      rule.exclude = [
        /node_modules\/@microsoft\/fetch-event-source/, // 忽略 fetch-event-source
        /node_modules\/@react-native/,                  // 忽略 React Native 相关模块
      ];
    }
    return rule;
  });

  return config;
};