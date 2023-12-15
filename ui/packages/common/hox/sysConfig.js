import { createGlobalStore } from 'hox';
import { useEffect, useState } from 'react';

import kfsConfig from '@kfs/common/kfsConfig/config.web';
import defaultConfig from '@kfs/common/kfsConfig/default';

function useFunc() {
  const [sysConfig, setSysConfig] = useState(kfsConfig.get() || defaultConfig);
  useEffect(() => {
    if (window.kfs.env.VITE_APP_PLATFORM === 'web' && window.kfs.env.MODE === 'production') {
      sysConfig.webServer = location.origin;
    }
    kfsConfig.set(sysConfig);
  }, [sysConfig]);
  return {
    sysConfig,
    setSysConfig,
    resetSysConfig: () => setSysConfig(defaultConfig),
  };
}

export const [useSysConfig, useSysConfig2] = createGlobalStore(useFunc);

export default useSysConfig;

export const getSysConfig = () => {
  return useSysConfig2().sysConfig;
}