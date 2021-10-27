import { useState, useEffect } from 'react';
import { createModel } from 'hox';
import { v4 } from 'uuid';

import kfsConfig from 'kfsConfig/index';

const defaultConfig = {
  clientID: v4(),
  theme: 'dark',
  backendProcess: {
    port: '1123',
    status: 'red',
  },
  remotes: [{
    name: '测试账号',
    type: '阿里云盘',
    loginType: 'refreshToken',
    refreshToken: '96246b97eb994fcaa4e8abb553d502bb',
  }],
  downloadPath: '',
};

function useSysConfig() {
  const initConfig = kfsConfig.get() || defaultConfig;
  window.clientID = initConfig.clientID;
  console.log(initConfig);
  const [sysConfig, setSysConfig] = useState(initConfig);
  useEffect(() => {
    kfsConfig.set(sysConfig);
  }, [sysConfig]);
  return {
    sysConfig,
    setSysConfig,
    resetSysConfig: () => setSysConfig(defaultConfig),
  };
}

export default createModel(useSysConfig);
