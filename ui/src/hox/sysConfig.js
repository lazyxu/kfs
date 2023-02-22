import { useState, useEffect } from 'react';
import { createGlobalStore } from 'hox';

import kfsConfig from 'kfsConfig/localStorage';
import defaultConfig from 'kfsConfig/default';

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

export const [useSysConfig, getSysConfig] = createGlobalStore(useFunc);

export default useSysConfig;
