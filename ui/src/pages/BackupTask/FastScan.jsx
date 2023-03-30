import {Alert, Button, Stack, TextField, Typography} from "@mui/material";
import BackupSizeStatus from "./BackupSizeStatus";
import {useState} from "react";
import useWebSocket from "react-use-websocket";
import {getSysConfig} from "../../hox/sysConfig";
import {v4 as uuid} from 'uuid';
import humanize from "humanize";

function isInvalidBackupDir(backupDir) {
    return backupDir === "";
}

let id;

export default function () {
    const sysConfig = getSysConfig().sysConfig;
    const {sendJsonMessage, lastJsonMessage} = useWebSocket("ws://127.0.0.1:" + sysConfig.port + "/ws");
    const [backupDir, setBackupDir] = useState('');
    const finished = lastJsonMessage?.finished;
    if (finished) {
        id = undefined;
    }
    console.log(id, lastJsonMessage)
    return (
        <Stack spacing={2}>
            <TextField variant="standard" label="本地文件夹路径" type="search" sx={{width: "50%"}}
                       value={backupDir}
                       onChange={e => setBackupDir(e.target.value)}/>
            {id ?
                <Button variant="outlined" sx={{width: "10em"}}
                        onClick={e => {
                            sendJsonMessage({type: "fastScan.cancel", id});
                        }}
                >
                    取消
                </Button>
                :
                <Button variant="outlined" sx={{width: "10em"}}
                        disabled={isInvalidBackupDir(backupDir)}
                        onClick={e => {
                            if (id) {
                                sendJsonMessage({type: "fastScan.cancel", id});
                            }
                            id = uuid();
                            console.log("fastScan", id, backupDir);
                            sendJsonMessage({type: "fastScan", id, data: {backupDir: backupDir}});
                        }}
                >
                    检测总大小
                </Button>}
            {lastJsonMessage && (lastJsonMessage.errMsg ?
                    <Alert variant="outlined" sx={{width: "max-content"}} severity="error">
                        {lastJsonMessage.errMsg}
                    </Alert>
                    :
                    <Alert variant="outlined" sx={{width: "max-content"}} severity={lastJsonMessage.finished ? "success" : "info"}>
                        <Typography>id：{lastJsonMessage.id}</Typography>
                        <Typography>待计算的文件和目录数量：{lastJsonMessage.data.stackSize}</Typography>
                        <Typography>文件数量：{lastJsonMessage.data.fileCount}</Typography>
                        <Typography>目录数量：{lastJsonMessage.data.dirCount}</Typography>
                        <Typography>总大小：{humanize.filesize(lastJsonMessage.data.fileSize)}</Typography>
                    </Alert>
            )}
        </Stack>
    );
}
