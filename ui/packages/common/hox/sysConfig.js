import { createGlobalStore } from 'hox';
import { useEffect, useState } from 'react';

const defaultConfig = {
  api: 'web',
  theme: 'dark',
  webServer: 'http://127.0.0.1:1123',
  socketServer: '127.0.0.1:1124',
  maxContentSize: 1 * 1024 * 1024,
  port: "11234",
};

function useFunc() {
  const [sysConfig, setSysConfig] = useState();
  useEffect(() => {
    (async () => {
      let c;
      if (window.getKfsConfig) {
        c = await window.getKfsConfig();
      } else if (window.kfsConfig) {
        c = window.kfsConfig;
      }
      if (!c) {
        c = defaultConfig;
      }
      if (window.kfsEnv.VITE_APP_PLATFORM === 'web' && window.kfsEnv.MODE === 'production') {
        c.webServer = location.origin;
      }
      console.log(c)
      setSysConfig(c);
    })();
  }, []);
  useEffect(() => {
    if (!sysConfig) {
      return;
    }
    if (window.setKfsConfig) {
      window.setKfsConfig(sysConfig); // async?
    } else {
      window.kfsConfig = sysConfig;
    }
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