import * as mockApi from "./mock/branch";
import * as webApi from "./web/branch";

export async function listBranch(setResourceManager) {
  console.log('api.listBranch');
  let branches = await webApi.listBranch();
  setResourceManager(prev => {
      return {...prev, branches, branchName: null, filePath: [], dirItems: null, file: null};
  });
}
