import {httpGet} from "./common";

export async function open(branchName, filePath) {
    console.log('web.open', branchName, filePath);
    return await httpGet("/api/v1/open", {
        branchName,
        filePath: filePath.join("/"),
        maxContentSize: 1024*1024,
    });
}
