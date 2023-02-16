import * as mockApi from "./mock/branch";

export async function listBranch(setResourceManager) {
  console.log('api.listBranch');
  let dirItems;
  let branches = await mockApi.listBranch();
  setResourceManager(prev => {
      return {...prev, branches, branchName: null, filePath: [], dirItems: null, file: null};
  });
}
