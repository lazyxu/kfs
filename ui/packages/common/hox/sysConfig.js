import { createGlobalStore } from 'hox';
import { useEffect, useState } from 'react';

import kfsConfig from '@kfs/common/kfsConfig/config.web';
import defaultConfig from '@kfs/common/kfsConfig/default';

function useFunc() {
  const [sysConfig, setSysConfig] = useState(kfsConfig.get() || defaultConfig);
  useEffect(() => {
    if (import.meta.env.REACT_APP_PLATFORM === 'web' && import.meta.env.NODE_ENV === 'production') {
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
