const createExpoWebpackConfigAsync = require('@expo/webpack-config');
const path = require('path');

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
  config.module.rules.push({
    test: /\.(js|jsx)$/,
    include: /(@kfs\/common)/, // 显式包含 workspace 包
    use: ['babel-loader']
  });
  config.module.rules.push({
      test: /\.(js|jsx)$/,
      // 包含node_modules中需要编译的特定库
      include: [
        // 添加下面这行，包含react-native-image-viewing的路径
        path.resolve(__dirname, '../node_modules/@staltz/react-native-image-viewing'),
      ],
      use: ['babel-loader']
    });
  return config;
};
