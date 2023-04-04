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

function isInvalidSrcPath(srcPath) {
    return srcPath === "";
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
    const [srcPath, setSrcPath] = useState('');
    const [dstPath, setDstPath] = useState('');
    const [finished, setFinished] = useState(true);
    useEffect(() => {
        if (!lastJsonMessage) {
            return;
        }
        setFinished(lastJsonMessage.finished);
    }, [lastJsonMessage]);
    return (
        <Stack spacing={2} style={{display: show ? undefined : "none"}}>
            <TextField variant="standard" label="本地文件夹路径" type="search" sx={{minWidth: "50%"}}
                       value={srcPath}
                       onChange={e => setSrcPath(e.target.value)}/>
            <TextField
                label="文件忽略规则"
                multiline minRows={4} maxRows={4}
                variant="standard"
                size="small"
            />
            <Stack spacing={2} direction="row">
                <FormControl sx={{minWidth: "10em"}} size="small">
                    <InputLabel id="demo-simple-select-label">备份分支</InputLabel>
                    <Select
                        labelId="demo-simple-select-label"
                        value={branchName}
                        onChange={e => setBranchName(e.target.value)}
                        autoWidth={true}
                    >
                        {branches.map(branch =>
                            <MenuItem key={branch.name} value={branch.name}>{branch.name}</MenuItem>
                        )}
                    </Select>
                </FormControl>
                <TextField variant="standard" label="远程文件夹路径" type="search" sx={{minWidth: "50%"}}
                           value={dstPath}
                           onChange={e => setDstPath(e.target.value)}
                           size="small"
                />
            </Stack>
            {!finished ?
                <Button variant="outlined" sx={{width: "10em"}}
                        onClick={e => {
                            sendJsonMessage({type: "cancel", id});
                        }}
                >
                    取消
                </Button>
                :
                <Button variant="outlined" sx={{width: "10em"}}
                        disabled={isInvalidSrcPath(srcPath) || isInvalidBranchName(branchName)}
                        onClick={e => {
                            let newId = uuid();
                            setId(newId);
                            console.log("fastBackup", newId, srcPath);
                            sendJsonMessage({
                                type: "fastBackup", id: newId, data: {
                                    srcPath,
                                    serverAddr: sysConfig.webServer,
                                    branchName,
                                    dstPath
                                }
                            });
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
                        <Typography>commitId：{lastJsonMessage.data.commitId}</Typography>
                        <Typography>文件数量：{lastJsonMessage.data.count}</Typography>
                        <Typography>总大小：{humanize.filesize(lastJsonMessage.data.size)}</Typography>
                    </Alert>
            ) : <Box/>}
        </Stack>
    );
}
