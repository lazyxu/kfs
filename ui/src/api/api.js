let mockApi = require("./mock/api");

export async function open(sysConfig, setResourceManager, branchName, filePath) {
    console.log('api.open', branchName, filePath, filePath.join('/'))
    let dirItems;
    let isDir = await mockApi.open(sysConfig, branchName, filePath, (content) => {
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
    console.log('api.list', branchName, filePath, filePath.join('/'))
    let dirItems;
    await mockApi.list(sysConfig, branchName, filePath, (total) => {
        dirItems = new Array(total);
    }, (dirItem, i) => {
        dirItems[i] = dirItem;
    });
    setResourceManager(prev => {
        return {...prev, branchName, filePath, dirItems, content: null};
    });
}
