import axios from 'axios';
import https from 'https';
import useNotification from 'common/components/Notification/Notification';

const httpsAgent = new https.Agent({
  ca: axios('./extraResources/rootCA.pem'),
  cert: axios('./extraResources/localhost.pem'),
  key: axios('./extraResources/localhost-key.pem'),
});

export function getBackendInstance() {
  return axios.create({
    baseURL: `https://localhost:${window.backendPort}`,
    httpsAgent,
  });
}

export async function post(url, data) {
  try {
    const res = await getBackendInstance().post(url, data);
    return res.data;
  } catch (e) {
    const msg = e?.response?.data?.message;
    if (msg) {
      throw new Error(msg);
    }
    console.error(e);
    throw e;
  }
}
