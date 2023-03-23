import {Button, Stack, TextField} from "@mui/material";
import BackupSizeStatus from "./BackupSizeStatus";
import {useState} from "react";
import useWebSocket from "react-use-websocket";
import {getSysConfig} from "../../hox/sysConfig";
import { v4 as uuid } from 'uuid';

function isInvalidBackupDir(backupDir) {
    return backupDir === "";
}

let id;

export default function () {
    const sysConfig = getSysConfig().sysConfig;
    const {sendJsonMessage, lastJsonMessage, readyState} = useWebSocket("ws://127.0.0.1:" + sysConfig.port + "/ws");
    const [backupDir, setBackupDir] = useState('');
    return (
        <Stack spacing={2}>
            <TextField variant="standard" label="本地文件夹路径" type="search" sx={{width: "50%"}}
                       value={backupDir}
                       onChange={e => setBackupDir(e.target.value)}/>
            <Button variant="outlined" sx={{width: "10em"}}
                    disabled={isInvalidBackupDir(backupDir)}
                    onClick={e => {
                        if (id) {
                            sendJsonMessage({type: "scan.cancel", id});
                        }
                        id = uuid();
                        console.log("scan", id, backupDir);
                        sendJsonMessage({type: "scan", id, data: {backupDir: backupDir}});
                    }}
            >
                检测总大小
            </Button>
            <Button variant="outlined" sx={{width: "10em"}}
                    disabled={isInvalidBackupDir(backupDir)}
                    onClick={e => {
                        id = uuid();
                        sendJsonMessage({type: "scan.cancel", id});
                    }}
            >
                取消
            </Button>
            <BackupSizeStatus json={lastJsonMessage}/>
        </Stack>
    );
}
