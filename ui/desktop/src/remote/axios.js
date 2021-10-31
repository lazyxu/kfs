import axios from 'axios';
import https from 'https';

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
