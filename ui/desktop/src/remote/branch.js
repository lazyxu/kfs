import { getBackendInstance } from './axios';

export async function listBranches() {
  console.log('listBranches ----------------');
  const res = await getBackendInstance().get('/api/branches');
  console.log('res is ----------------', res);
  return res.data;
}
