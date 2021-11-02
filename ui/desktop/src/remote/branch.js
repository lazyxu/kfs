import { getBackendInstance } from './axios';

export async function listBranches() {
  const res = await getBackendInstance().get('/api/listBranches');
  return res.data;
}

export async function createBranch(clientID, branchName) {
  const res = await getBackendInstance().post('/api/createBranch', {
    clientID, branchName,
  });
  return res.data;
}

export async function deleteBranch(clientID, branchName) {
  const res = await getBackendInstance().post('/api/deleteBranch', {
    clientID, branchName,
  });
  return res.data;
}

export async function renameBranch(clientID, old, _new) {
  const res = await getBackendInstance().post('/api/renameBranch', {
    clientID, old, new: _new,
  });
  return res.data;
}
