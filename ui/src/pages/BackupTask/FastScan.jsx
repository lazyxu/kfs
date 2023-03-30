import {Alert, Box, Button, Stack, TextField, Typography} from "@mui/material";
import {useEffect, useState} from "react";
import useWebSocket from "react-use-websocket";
import {getSysConfig} from "../../hox/sysConfig";
import {v4 as uuid} from 'uuid';
import humanize from "humanize";

function isInvalidBackupDir(backupDir) {
    return backupDir === "";
}

export default function ({show}) {
    const sysConfig = getSysConfig().sysConfig;
    const [id, setId] = useState();
    const {sendJsonMessage, lastJsonMessage} = useWebSocket("ws://127.0.0.1:" + sysConfig.port + "/ws", {
        share: true,
        filter: message => {
            if (!(message?.data)) {
                return false;
            }
            let curId = JSON.parse(message.data)?.id;
            if (curId !== id) {
                return false;
            }
            return true;
        }
    });
    const [backupDir, setBackupDir] = useState('');
    const [finished, setFinished] = useState(true);
    useEffect(()=> {
        if (!lastJsonMessage) {
            return;
        }
        setFinished(lastJsonMessage.finished);
    }, [lastJsonMessage]);
    if (!show) {
        return
    }
    return (
        <Stack spacing={2}>
            <TextField variant="standard" label="本地文件夹路径" type="search" sx={{width: "50%"}}
                       value={backupDir}
                       onChange={e => setBackupDir(e.target.value)}/>
            {!finished ?
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
                            let newId = uuid();
                            setId(newId);
                            console.log("fastScan", newId, backupDir);
                            sendJsonMessage({type: "fastScan", id: newId, data: {backupDir: backupDir}});
                        }}
                >
                    快速扫描
                </Button>}
            {lastJsonMessage ? (lastJsonMessage.errMsg ?
                    <Alert variant="outlined" sx={{width: "max-content"}} severity="error">
                        {lastJsonMessage.errMsg}
                    </Alert>
                    :
                    <Alert variant="outlined" sx={{width: "max-content"}}
                           severity={lastJsonMessage.finished ? "success" : "info"}>
                        <Typography>id：{id}</Typography>
                        <Typography>待计算的文件和目录数量：{lastJsonMessage.data.stackSize}</Typography>
                        <Typography>文件数量：{lastJsonMessage.data.fileCount}</Typography>
                        <Typography>目录数量：{lastJsonMessage.data.dirCount}</Typography>
                        <Typography>总大小：{humanize.filesize(lastJsonMessage.data.fileSize)}</Typography>
                    </Alert>
            ) : <Box/>}
        </Stack>
    );
}
