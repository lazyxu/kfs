import { post } from './axios';

export async function getBranchHash(clientID, branchName) {
  return post('/api/getBranchHash', {
    clientID, branchName,
  });
}

export async function listBranches() {
  return post('/api/listBranches');
}

export async function createBranch(clientID, branchName) {
  return post('/api/createBranch', {
    clientID, branchName,
  });
}

export async function deleteBranch(clientID, branchName) {
  return post('/api/deleteBranch', {
    clientID, branchName,
  });
}

export async function renameBranch(clientID, old, _new) {
  return post('/api/renameBranch', {
    clientID, old, new: _new,
  });
}
