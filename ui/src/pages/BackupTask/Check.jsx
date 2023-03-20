import {Button, FormControl, InputLabel, MenuItem, Select, Stack, TextField, Typography} from "@mui/material";
import {useEffect, useState} from "react";
import {getBranchApi} from "../../api/branch";
import useWebSocket, {ReadyState} from "react-use-websocket";
import useSysConfig from "../../hox/sysConfig";
import BackupSizeStatus from "./BackupSizeStatus";

function isInvalidBackupDir(backupDir) {
    return backupDir === "";
}

function isInvalidBranchName(branchName) {
    return branchName === "";
}

let lastId = 0;

export default function () {
    const {sysConfig, setSysConfig, resetSysConfig} = useSysConfig();
    const {sendJsonMessage, lastJsonMessage, readyState} = useWebSocket("ws://127.0.0.1:" + sysConfig.port + "/ws");
    const [id, setId] = useState(lastId);
    const [branches, setBranches] = useState([]);
    useEffect(() => {
        getBranchApi().listBranch().then(setBranches);
    }, []);
    const [branchName, setBranchName] = useState('');
    const [backupDir, setBackupDir] = useState('');
    const [calculateBackupSizeResult, setCalculateBackupSizeResult] = useState('');
    const [backupResult, setBackupResult] = useState('');
    const [page, setPage] = useState(0);
    return (
        <Stack spacing={2}>
            <TextField variant="standard" label="本地文件夹路径" type="search" sx={{width: "100%"}}
                       value={backupDir}
                       onChange={e => setBackupDir(e.target.value)}/>
            <Button variant="outlined" sx={{width: "10em"}}
                    disabled={isInvalidBackupDir(backupDir)}
                    onClick={e => {
                        sendJsonMessage({type: "calculateBackupSize.cancel", id, data: {backupDir: backupDir}});
                        const newId = id + 1;
                        setId(newId);
                        console.log("calculateBackupSize", newId, backupDir);
                        sendJsonMessage({type: "calculateBackupSize", id: newId, data: {backupDir: backupDir}});
                    }}
            >
                检测总大小
            </Button>
            <Button variant="outlined" sx={{width: "10em"}}
                    disabled={isInvalidBackupDir(backupDir)}
                    onClick={e => {
                        sendJsonMessage({type: "calculateBackupSize.cancel", id, data: {backupDir: backupDir}});
                    }}
            >
                取消
            </Button>
            <BackupSizeStatus json={lastJsonMessage}/>
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
            <Button variant="outlined" sx={{width: "10em"}}
                    disabled={isInvalidBackupDir(backupDir) || isInvalidBranchName(branchName)}
                    onClick={e => {
                        console.log("backup", backupDir, branchName)
                    }}
            >
                开始备份
            </Button>
            <Typography>{backupResult}</Typography>
        </Stack>
    );
}
