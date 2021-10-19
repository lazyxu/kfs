import { useState, useEffect } from 'react';
import { createModel } from 'hox';

import kfsConfig from 'common/kfsConfig/localStorage';

const defaultConfig = {
  theme: 'dark',
  username: '17161951517',
  refreshToken: '96246b97eb994fcaa4e8abb553d502bb',
};

function useSysConfig() {
  const [sysConfig, setSysConfig] = useState(kfsConfig.get() || defaultConfig);
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
