let mockApi = require("./mock/api");

export async function open(sysConfig, setResourceManager, branchName, filePath) {
    console.log('api.open', branchName, filePath, filePath.join('/'));
    let dirItems;
    let isDir = await mockApi.open(branchName, filePath, (content) => {
        setResourceManager(prev => {
            return {
                ...prev, branchName, filePath,
                dirItems: null,
                content,
            };
        });
    }, (total) => {
        dirItems = new Array(total);
    }, (dirItem, i) => {
        dirItems[i] = dirItem;
    });
    if (isDir) {
        setResourceManager(prev => {
            return {
                ...prev, branchName, filePath,
                dirItems: dirItems ? dirItems : prev.dirItems,
                content: null,
            };
        });
    }
}

export async function list(sysConfig, setResourceManager, branchName, filePath) {
    console.log('api.list', branchName, filePath, filePath.join('/'));
    let dirItems;
    await mockApi.list(branchName, filePath, (total) => {
        dirItems = new Array(total);
    }, (dirItem, i) => {
        dirItems[i] = dirItem;
    });
    setResourceManager(prev => {
        return {...prev, branchName, filePath, dirItems, content: null};
    });
}

export async function newFile(sysConfig, setResourceManager, branchName, filePath) {
    console.log('api.newFile', branchName, filePath, filePath.join('/'));
    await mockApi.newFile(branchName, filePath);
    await list(sysConfig, setResourceManager, branchName, filePath)
}

export async function newDir(sysConfig, setResourceManager, branchName, filePath) {
    console.log('api.newDir', branchName, filePath, filePath.join('/'));
    await mockApi.newDir(branchName, filePath);
    await list(sysConfig, setResourceManager, branchName, filePath)
}
