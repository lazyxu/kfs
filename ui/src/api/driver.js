import { getSysConfig } from "../hox/sysConfig";
import * as mockApi from "./mock/driver";
import * as webApi from "./web/driver";

export function getDriverApi() {
    return getSysConfig().sysConfig.api === "web" ? webApi : mockApi;
}

export async function listDriver(setResourceManager) {
    console.log('api.listDriver');
    let drivers = await getDriverApi().listDriver();
    setResourceManager(prev => {
        return {...prev, drivers, driverName: null, filePath: [], dirItems: null, file: null};
    });
}

export async function newDriver(setResourceManager, name, description, type, code) {
    console.log('api.newDriver', name, description, type, code);
    // TODO: exist
    await getDriverApi().newDriver(name, description, type, code);
    await listDriver(setResourceManager);
}

export async function deleteDriver(setResourceManager, name) {
    console.log('api.deleteDriver', name);
    await getDriverApi().deleteDriver(name);
    await listDriver(setResourceManager);
}
