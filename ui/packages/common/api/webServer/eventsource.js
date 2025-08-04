import { fetchEventSource } from "@kfs/common/api/webServer/fetchEventSource";
import { getSysConfig } from "../../hox/sysConfig";

export function listDCIM(signal, onopen, onJson, onerror) {
        fetchEventSource(`${getSysConfig().webServer}/api/v1/event/listDCIM`, {
            signal,
            openWhenHidden: false,
            onopen(response) {
                if (response.ok && response.headers.get('content-type').includes("text/event-stream")) {
                    return; // everything's good
                }
                console.error('onopen', response);
                onerror?.('onopen', response);
                onopen?.(response);
            },
            onmessage(msg) {
                if (msg.event === 'FatalError') {
                    console.error('onmessage', msg);
                    onerror?.('onmessage', msg);
                    return;
                }
                let json = JSON.parse(msg.data);
                console.log('listDCIM.onmessage:', json);
                if (json.errMsg) {
                    console.error('listDCIM.onmessage.json', json.errMsg);
                    onerror?.('listDCIM.onmessage.json', json.errMsg);
                    return;
                }
                onJson?.(json);
            },
            onclose() {
                console.error('close');
            },
            onerror(err) {
                console.error('err', err);
                onerror?.('onerror', err);
            }
        });
}
