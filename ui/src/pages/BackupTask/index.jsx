import {Box, Stack, Tab, Tabs} from "@mui/material";
import {useState} from "react";
import Scan from "./Scan";
import FastBackup from "./FastBackup";

export default function ({show}) {
    const [page, setPage] = useState(0);
    return (
        <Stack spacing={2} sx={{
            position: "absolute", width: "100%", height: "100%"
        }} style={{display: show ? undefined : "none"}}>
            <Box sx={{borderBottom: 1, borderColor: 'divider', width: "100%"}}>
                <Tabs value={page} variant="scrollable" scrollButtons="auto" onChange={(e, v) => setPage(v)}>
                    <Tab label="主页"/>
                    <Tab label="扫描"/>
                    <Tab label="扫描历史"/>
                    <Tab label="快速备份"/>
                    <Tab label="备份"/>
                    <Tab label="备份历史"/>
                </Tabs>
            </Box>
            <Box sx={{padding: "1em"}}>
                <Scan show={page === 1}/>
                <FastBackup show={page === 0}/>
            </Box>
        </Stack>
    );
}
