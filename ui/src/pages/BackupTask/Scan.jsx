import {
    Alert,
    Box,
    Button,
    Checkbox,
    FormControlLabel,
    FormGroup,
    LinearProgress,
    Stack,
    TextField,
    Typography
} from "@mui/material";
import {useEffect, useState} from "react";
import useWebSocket from "react-use-websocket";
import {getSysConfig} from "../../hox/sysConfig";
import {v4 as uuid} from 'uuid';
import humanize from "humanize";
import moment from "moment/moment";
import LinearProgressWithLabel from "./LinearProgressWithLabel";

function isInvalidSrcPath(srcPath) {
    return srcPath === "";
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
    const [srcPath, setSrcPath] = useState('');
    const [finished, setFinished] = useState(true);
    const [timeCost, setTimeCost] = useState(null);
    useEffect(() => {
        if (!lastJsonMessage) {
            return;
        }
        setTimeCost(prev => {
            if (lastJsonMessage !== prev) {
                return moment().diff(startTime, 'seconds', true);
            }
            return null;
        });
        setFinished(lastJsonMessage.finished);
    }, [lastJsonMessage]);
    const [startTime, setStartTime] = useState();
    const [record, setRecord] = useState(false);
    const [concurrent, setConcurrent] = useState(15);
    return (
        <Stack spacing={2} style={{display: show ? undefined : "none"}}>
            <TextField variant="standard" label="本地文件夹路径" type="search" sx={{width: "50%"}}
                       value={srcPath}
                       onChange={e => setSrcPath(e.target.value)}/>
            <TextField
                label="文件忽略规则"
                multiline minRows={4} maxRows={4}
                variant="standard"
                size="small"
            />
            <FormGroup>
                <FormControlLabel control={<Checkbox value={record} onChange={e => {
                    setRecord(e.target.checked)
                }}/>} label="记录文件大小"/>
                <TextField
                    label="并发扫描数量"
                    variant="standard"
                    size="small"
                    inputProps={{inputMode: 'numeric', pattern: "[1-9][0-9]*"}}
                    value={concurrent}
                    onChange={e => {
                        const val = e.target.value;
                        console.log(val, val.match(/^[1-9][0-9]*$/))
                        if (!val.match(/^[1-9][0-9]*$/)) {
                            return e.preventDefault();
                        }
                        setConcurrent(parseInt(val, 10));
                    }}/>
            </FormGroup>
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
                        disabled={isInvalidSrcPath(srcPath)}
                        onClick={e => {
                            let newId = uuid();
                            setId(newId);
                            let data = {record, srcPath, concurrent};
                            console.log("scan", newId, data);
                            setStartTime(moment());
                            sendJsonMessage({type: "scan", id: newId, data});
                        }}
                >
                    扫描
                </Button>
            }
            {lastJsonMessage ? (lastJsonMessage.errMsg ?
                    <Alert variant="outlined" sx={{width: "max-content"}} severity="error">
                        {lastJsonMessage.errMsg}
                    </Alert>
                    :
                    <Alert variant="outlined" sx={{width: "max-content"}}
                           severity={lastJsonMessage.finished ? "success" : "info"}>
                        <Typography>耗时：{timeCost} 秒</Typography>
                        <LinearProgressWithLabel variant="determinate" value={lastJsonMessage.data.fileCount/(lastJsonMessage.data.fileCount+lastJsonMessage.data.stackSize)*100} />
                        <Typography>待计算的文件和目录数量：{lastJsonMessage.data.stackSize}</Typography>
                        <Typography>文件数量：{lastJsonMessage.data.fileCount}</Typography>
                        <Typography>目录数量：{lastJsonMessage.data.dirCount}</Typography>
                        <Typography>总大小：{humanize.filesize(lastJsonMessage.data.fileSize)}</Typography>
                    </Alert>
            ) : <Box/>}
        </Stack>
    );
}
