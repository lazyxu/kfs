import {Box, Stack, Tab, Tabs, Typography} from "@mui/material";
import {useState} from "react";
import useWebSocket, {ReadyState} from "react-use-websocket";
import useSysConfig from "../../hox/sysConfig";
import FastScan from "./FastScan";
import Scan from "./Scan";
import FastBackup from "./FastBackup";

export default function ({show}) {
    const [page, setPage] = useState(0);
    if (!show) {
        return;
    }
    return (
        <Stack spacing={2} sx={{
            position: "absolute", width: "100%", height: "100%"
        }}>
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
                <FastScan show={page === 1}/>
                <Scan show={page === 2}/>
                <FastBackup show={page === 4}/>
            </Box>
        </Stack>
    );
}
