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
import FastScan from "./FastScan";
import Scan from "./Scan";
import FastBackup from "./FastBackup";

export default function () {
    const {sysConfig, setSysConfig, resetSysConfig} = useSysConfig();
    const {sendJsonMessage, lastJsonMessage, readyState} = useWebSocket("ws://127.0.0.1:" + sysConfig.port + "/ws", {share: true});
    const connectionStatus = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];
    const [page, setPage] = useState(0);
    return (
        <Stack spacing={2} sx={{
            position: "absolute", width: "100%", height: "100%"
        }}>
            <Typography>The WebSocket is currently {connectionStatus}</Typography>
            <Box sx={{borderBottom: 1, borderColor: 'divider', width: "100%"}}>
                <Tabs value={page} variant="scrollable" scrollButtons="auto" onChange={(e, v) => setPage(v)}>
                    <Tab label="主页"/>
                    <Tab label="快速扫描"/>
                    <Tab label="扫描"/>
                    <Tab label="扫描历史"/>
                    <Tab label="快速备份"/>
                    <Tab label="备份"/>
                    <Tab label="备份历史"/>
                </Tabs>
            </Box>
            <Box sx={{padding: "1em"}}>
                {page === 1 && <FastScan/>}
                {page === 2 && <Scan/>}
                {page === 4 && <FastBackup/>}
            </Box>
        </Stack>
    );
}
