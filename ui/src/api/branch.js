import * as mockApi from "./mock/branch";
import * as webApi from "./web/branch";
import {getSysConfig} from "../hox/sysConfig";

function getBranchApi() {
    return getSysConfig().sysConfig.api === "web" ? webApi : mockApi;
}

export async function listBranch(setResourceManager) {
    console.log('api.listBranch');
    let branches = await getBranchApi().listBranch();
    setResourceManager(prev => {
        return {...prev, branches, branchName: null, filePath: [], dirItems: null, file: null};
    });
}
