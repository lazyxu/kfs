import { useState, useEffect } from 'react';
import { createGlobalStore } from 'hox';

import kfsConfig from 'kfsConfig/localStorage';

const defaultConfig = {
  theme: 'dark',
  username: '17161951517',
  refreshToken: '96246b97eb994fcaa4e8abb553d502bb',
};

function useFunc() {
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

const [useSysConfig] = createGlobalStore(useFunc);

export default useSysConfig;
