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
    FormControlLabel,
    Switch,
    Typography
} from "@mui/material";
import { useEffect, useState } from "react";
import useWebSocket from "react-use-websocket";
import { getSysConfig } from "../../hox/sysConfig";
import { v4 as uuid } from 'uuid';
import humanize from "humanize";
import { getBranchApi } from "../../api/branch";
import AsyncSelect from "components/AsyncSelect";

function isInvalidSrcPath(srcPath) {
    return srcPath === "";
}

function isInvalidBranchName(branchName) {
    return branchName === "";
}

export default function ({ show }) {
    const sysConfig = getSysConfig().sysConfig;
    const [id, setId] = useState();
    const { sendJsonMessage, lastJsonMessage } = useWebSocket("ws://127.0.0.1:" + sysConfig.port + "/ws", {
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
    const [concurrent, setConcurrent] = useState(2);
    const [encoder, setEncoder] = useState("none");
    const [verbose, setVerbose] = useState(true);
    const [srcPath, setSrcPath] = useState('');
    const [dstPath, setDstPath] = useState('');
    const [errs, setErrs] = useState([]);
    const [finished, setFinished] = useState(true);
    useEffect(() => {
        if (!lastJsonMessage) {
            return;
        }
        setFinished(lastJsonMessage.finished);
        if (!lastJsonMessage.data) {
            return;
        }
        console.log(lastJsonMessage)
        if (lastJsonMessage.data.err) {
            setErrs(prev => prev.push({ err: lastJsonMessage.data.err, filePath: lastJsonMessage.data.filePath }));
        }
    }, [lastJsonMessage]);
    return (
        <Stack spacing={2} style={{ display: show ? undefined : "none" }}>
            <TextField variant="standard" label="本地文件夹路径" type="search" sx={{ minWidth: "50%" }}
                value={srcPath}
                onChange={e => setSrcPath(e.target.value)} />
            <Stack spacing={2} direction="row">
                <FormControl sx={{ minWidth: "10em" }}>
                    <AsyncSelect
                        label="备份分支"
                        fetchOptions={async () => {
                            let branches = await getBranchApi().listBranch();
                            return branches.map(branch => branch.name);
                        }}
                        value={branchName}
                        onChange={name => setBranchName(name)}
                    />
                </FormControl>
                <TextField variant="standard" label="远程文件夹路径" type="search" sx={{ minWidth: "50%" }}
                    value={dstPath}
                    onChange={e => setDstPath(e.target.value)}
                />
            </Stack>
            <FormControl sx={{ minWidth: "10em" }}>
                <InputLabel id="backup-encoder-label">上传时压缩</InputLabel>
                <Select
                    labelId="backup-encoder-label"
                    value={encoder}
                    onChange={e => setEncoder(e.target.value)}
                    sx={{ width: "10em" }}
                >
                    {["none", "lz4"].map(value =>
                        <MenuItem key={value} value={value}>{value}</MenuItem>
                    )}
                </Select>
            </FormControl>
            <FormControl sx={{ minWidth: "10em" }}>
                <InputLabel id="backup-concurrent-label">同时上传文件数</InputLabel>
                <Select
                    labelId="backup-concurrent-label"
                    value={concurrent}
                    onChange={e => setConcurrent(e.target.value)}
                    sx={{ width: "10em" }}
                >
                    {[1, 2, 3, 4, 5].map(value =>
                        <MenuItem key={value} value={value}>{value}</MenuItem>
                    )}
                </Select>
            </FormControl>
            <FormControlLabel label="显示上传进度" control={
                <Switch
                    checked={verbose}
                    onChange={e => setVerbose(e.target.checked)}
                />
            } />
            {!finished ?
                <Button variant="outlined" sx={{ width: "10em" }}
                    onClick={e => {
                        sendJsonMessage({ type: "cancel", id });
                    }}
                >
                    取消
                </Button>
                :
                <Button variant="outlined" sx={{ width: "10em" }}
                    disabled={isInvalidSrcPath(srcPath) || isInvalidBranchName(branchName)}
                    onClick={e => {
                        let newId = uuid();
                        setId(newId);
                        setErrs([]);
                        console.log("fastBackup", newId, srcPath);
                        sendJsonMessage({
                            type: "fastBackup", id: newId, data: {
                                srcPath, verbose, concurrent, encoder,
                                serverAddr: sysConfig.webServer,
                                branchName,
                                dstPath
                            }
                        });
                    }}
                >
                    快速备份
                </Button>}
            {lastJsonMessage?.errMsg &&
                <Alert variant="outlined" sx={{ width: "max-content" }} severity="error">
                    {lastJsonMessage.errMsg}
                </Alert>}
            {lastJsonMessage?.finished && lastJsonMessage?.data?.branch &&
                <Alert variant="outlined" sx={{ width: "max-content" }} severity={"success"}>
                    <Typography>id：{lastJsonMessage.id}</Typography>
                    <Typography>commitId：{lastJsonMessage.data.branch.commitId}</Typography>
                    <Typography>文件数量：{lastJsonMessage.data.branch.count}</Typography>
                    <Typography>总大小：{humanize.filesize(lastJsonMessage.data.branch.size)}</Typography>
                </Alert>}
            {lastJsonMessage?.data?.filePath && <Alert variant="outlined" sx={{ width: "max-content" }} severity={"info"}>
                <Typography>备份中：{humanize.filesize(lastJsonMessage.data.size)}/{humanize.filesize(lastJsonMessage.data.totalSize)}</Typography>
                <Typography>文件：{lastJsonMessage.data.fileCount}/{lastJsonMessage.data.totalFileCount}</Typography>
                <Typography>目录：{lastJsonMessage.data.dirCount}/{lastJsonMessage.data.totalDirCount}</Typography>
                <Typography>{lastJsonMessage.data.exist ? "已经存在" : "上传成功"}：{lastJsonMessage.data.filePath}</Typography>
            </Alert>}
            {errs.map(err =>
                <Alert variant="outlined" sx={{ width: "max-content" }} severity="error" key={err.filePath}>
                    {err.filePath}: {err.err}
                </Alert>
            )}
        </Stack>
    );
}
