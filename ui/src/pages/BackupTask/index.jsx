import {
    Box,
    Button,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    Stack,
    Tab,
    Tabs,
    TextField,
    Typography
} from "@mui/material";
import {useEffect, useState} from "react";
import {getBranchApi} from "../../api/branch";
import useWebSocket, {ReadyState} from "react-use-websocket";
import useSysConfig from "../../hox/sysConfig";
import BackupSizeStatus from "./BackupSizeStatus";
import FastCheck from "./FastCheck";

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
    const connectionStatus = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];
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
        <Stack spacing={2} sx={{
            position: "absolute", width: "100%", height: "100%"
        }}>
            <Typography>The WebSocket is currently {connectionStatus}</Typography>
            <Box sx={{borderBottom: 1, borderColor: 'divider', width: "100%"}}>
                <Tabs value={page} variant="scrollable" scrollButtons="auto" onChange={(e, v) => setPage(v)}>
                    <Tab label="主页"/>
                    <Tab label="历史"/>
                    <Tab label="快速检测"/>
                    <Tab label="带记录的检测"/>
                    <Tab label="备份"/>
                    <Tab label="其它"/>
                    <Tab label="其它"/>
                    <Tab label="其它"/>
                    <Tab label="其它"/>
                    <Tab label="其它"/>
                    <Tab label="其它"/>
                    <Tab label="其它"/>
                    <Tab label="其它3"/>
                    <Tab label="其它2"/>
                    <Tab label="其它1"/>
                </Tabs>
            </Box>
            <Box sx={{padding: "1em"}}>
                {page === 2 && <FastCheck/>}
                <Box sx={{borderBottom: 1, borderColor: 'divider', width: "100%"}}/>
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
            </Box>
        </Stack>
    );
}
