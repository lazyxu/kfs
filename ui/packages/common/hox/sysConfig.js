import { createGlobalStore } from 'hox';
import { useEffect, useState } from 'react';

import kfsConfig from '@kfs/common/kfsConfig/config.web';
import defaultConfig from '@kfs/common/kfsConfig/default';
import { getEnv } from './env';

function useFunc() {
  const [sysConfig, setSysConfig] = useState(kfsConfig.get() || defaultConfig);
  useEffect(() => {
    if (getEnv().VITE_APP_PLATFORM === 'web' && getEnv().MODE === 'production') {
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

export const [useSysConfig, getSysConfig] = createGlobalStore(useFunc);

export default useSysConfig;
