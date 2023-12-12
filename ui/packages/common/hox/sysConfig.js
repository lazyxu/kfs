import { createGlobalStore } from 'hox';
import { useEffect, useState } from 'react';

import defaultConfig from 'kfsConfig/default';
import kfsConfig from 'kfsConfig/localStorage';

function useFunc() {
  const [sysConfig, setSysConfig] = useState(kfsConfig.get() || defaultConfig);
  useEffect(() => {
    if (process.env.REACT_APP_PLATFORM === 'web' && process.env.NODE_ENV === 'production') {
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
