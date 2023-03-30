import {
    Alert,
    Box,
    Button,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    Stack,
    TextField,
    Typography
} from "@mui/material";
import {useEffect, useState} from "react";
import useWebSocket from "react-use-websocket";
import {getSysConfig} from "../../hox/sysConfig";
import {v4 as uuid} from 'uuid';
import humanize from "humanize";
import {getBranchApi} from "../../api/branch";

function isInvalidBackupDir(backupDir) {
    return backupDir === "";
}

function isInvalidBranchName(branchName) {
    return branchName === "";
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
    const [branches, setBranches] = useState([]);
    useEffect(() => {
        getBranchApi().listBranch().then(setBranches);
    }, []);
    const [branchName, setBranchName] = useState('');
    const [backupDir, setBackupDir] = useState('');
    const [finished, setFinished] = useState(true);
    useEffect(()=> {
        if (!lastJsonMessage) {
            return;
        }
        setFinished(lastJsonMessage.finished);
    }, [lastJsonMessage]);
    return (
        <Stack spacing={2} style={{display: show?undefined:"none"}}>
            <FormControl sx={{width: "10em"}}>
                <InputLabel id="demo-simple-select-label">备份分支</InputLabel>
                <Select
                    labelId="demo-simple-select-label"
                    value={branchName}
                    onChange={e => setBranchName(e.target.value)}
                >
                    {branches.map(branch =>
                        <MenuItem key={branch.name} value={branch.name}>{branch.name}</MenuItem>
                    )}
                </Select>
            </FormControl>
            <TextField variant="standard" label="本地文件夹路径" type="search" sx={{width: "50%"}}
                       value={backupDir}
                       onChange={e => setBackupDir(e.target.value)}/>
            {!finished ?
                <Button variant="outlined" sx={{width: "10em"}}
                        onClick={e => {
                            sendJsonMessage({type: "fastBackup.cancel", id});
                        }}
                >
                    取消
                </Button>
                :
                <Button variant="outlined" sx={{width: "10em"}}
                        disabled={isInvalidBackupDir(backupDir) || isInvalidBranchName(branchName)}
                        onClick={e => {
                            if (id) {
                                sendJsonMessage({type: "fastBackup.cancel", id});
                            }
                            let newId = uuid();
                            setId(newId);
                            console.log("fastBackup", newId, backupDir);
                            sendJsonMessage({type: "fastBackup", id: newId, data: {backupDir: backupDir}});
                        }}
                >
                    快速备份
                </Button>}
            {lastJsonMessage ? (lastJsonMessage.errMsg ?
                    <Alert variant="outlined" sx={{width: "max-content"}} severity="error">
                        {lastJsonMessage.errMsg}
                    </Alert>
                    :
                    <Alert variant="outlined" sx={{width: "max-content"}}
                           severity={lastJsonMessage.finished ? "success" : "info"}>
                        <Typography>id：{lastJsonMessage.id}</Typography>
                        <Typography>待计算的文件和目录数量：{lastJsonMessage.data.stackSize}</Typography>
                        <Typography>文件数量：{lastJsonMessage.data.fileCount}</Typography>
                        <Typography>目录数量：{lastJsonMessage.data.dirCount}</Typography>
                        <Typography>总大小：{humanize.filesize(lastJsonMessage.data.fileSize)}</Typography>
                    </Alert>
            ) : <Box/>}
        </Stack>
    );
}
