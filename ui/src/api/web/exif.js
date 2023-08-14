import {httpGet} from "./common";
import axios from "axios";
import {getSysConfig} from "hox/sysConfig";

export async function analysisExif(start) {
    console.log('web.analysisExif', start);
    return await httpGet("/api/v1/analysisExif", {
        start,
    });
}
