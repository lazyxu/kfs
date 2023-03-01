import * as mockApi from "./mock/branch";
import * as webApi from "./web/branch";
import {getSysConfig} from "../hox/sysConfig";

export function getBranchApi() {
    return getSysConfig().sysConfig.api === "web" ? webApi : mockApi;
}

export async function listBranch(setResourceManager) {
    console.log('api.listBranch');
    let branches = await getBranchApi().listBranch();
    setResourceManager(prev => {
        return {...prev, branches, branchName: null, filePath: [], dirItems: null, file: null};
    });
}

export async function newBranch(setResourceManager, name) {
    console.log('api.newBranch', name);
    // TODO: exist
    await getBranchApi().newBranch(name);
    await listBranch(setResourceManager);
}

export async function deleteBranch(setResourceManager, name) {
    console.log('api.deleteBranch', name);
    await getBranchApi().deleteBranch(name);
    await listBranch(setResourceManager);
}
