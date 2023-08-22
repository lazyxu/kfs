import { Box, Button, FormControlLabel, InputLabel, Radio, RadioGroup, Stack, TextField, Typography } from '@mui/material';
import { toPrecent } from 'api/utils/api';
import { getDiskUsage } from 'api/web/disk';
import useSysConfig from 'hox/sysConfig';
import humanize from 'humanize';
import { useEffect, useState } from 'react';

export default ({ show }) => {
    const { sysConfig, setSysConfig, resetSysConfig } = useSysConfig();
    const [diskUsage, setDiskUsage] = useState();
    useEffect(() => {
        getDiskUsage().then(setDiskUsage);
    }, []);
    return (
        <Stack style={{ padding: "1em", display: show ? undefined : "none" }}>
            <Box style={{ textAlign: "left", fontSize: "20px", paddingTop: "10px", paddingBottom: "10px" }}>储存空间</Box>
            {!diskUsage ? <Typography>正在获取...</Typography> :
                <Box>
                    <Box sx={{ width: "100%", height: "1em", paddingBottom: "1em" }}>
                        <Box sx={{ borderRadius: "5px 0 0 5px", display: "inline-block", height: "1em", width: toPrecent((diskUsage.total - diskUsage.file - diskUsage.metadata - diskUsage.thumbnail - diskUsage.free) / diskUsage.total), background: "gray" }}></Box>
                        <Box sx={{ display: "inline-block", height: "1em", width: toPrecent((diskUsage.file) / diskUsage.total), background: "green" }}></Box>
                        <Box sx={{ display: "inline-block", height: "1em", width: toPrecent((diskUsage.metadata) / diskUsage.total), background: "orange" }}></Box>
                        <Box sx={{ display: "inline-block", height: "1em", width: toPrecent((diskUsage.thumbnail) / diskUsage.total), background: "yellow" }}></Box>
                        <Box sx={{ borderRadius: "0 5px 5px 0", display: "inline-block", height: "1em", width: toPrecent((diskUsage.free) / diskUsage.total), background: "white" }}></Box>
                    </Box>
                    <Box>
                        <InputLabel>硬盘总空间：{humanize.filesize(diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "gray", height: "1em", width: "1em", borderRadius: "1em" }}></Box> 其他：{humanize.filesize(diskUsage.total - diskUsage.file - diskUsage.metadata - diskUsage.thumbnail - diskUsage.free)} {toPrecent((diskUsage.total - diskUsage.file - diskUsage.metadata - diskUsage.thumbnail - diskUsage.free) / diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "green", height: "1em", width: "1em", borderRadius: "1em" }}></Box> 云盘文件：{humanize.filesize(diskUsage.file)} {toPrecent(diskUsage.file / diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "orange", height: "1em", width: "1em", borderRadius: "1em" }}></Box> 元数据：{humanize.filesize(diskUsage.metadata)} {toPrecent(diskUsage.metadata / diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "yellow", height: "1em", width: "1em", borderRadius: "1em" }}></Box> 缩略图：{humanize.filesize(diskUsage.thumbnail)} {toPrecent(diskUsage.thumbnail / diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "white", height: "1em", width: "1em", borderRadius: "1em" }}></Box> 剩余：{humanize.filesize(diskUsage.free)} {toPrecent(diskUsage.free / diskUsage.total)}</InputLabel>
                    </Box>
                </Box>}
            <Box style={{ textAlign: "left", fontSize: "20px", paddingTop: "10px", paddingBottom: "10px" }}>设置</Box>
            {!sysConfig ? <span>加载中...</span>
                : (
                    <>
                        <Button variant="outlined" sx={{ width: "10em" }} onClick={e => resetSysConfig()} >恢复默认设置</Button>
                        <Box>
                            <InputLabel sx={{ display: "inline" }}>主题：</InputLabel>
                            <RadioGroup sx={{ display: "inline" }}
                                row
                                value={sysConfig.theme}
                                onChange={e => setSysConfig(c => ({ ...c, theme: e.target.value }))}
                                size="small"
                            >
                                {["light", "dark", "system"].map(value =>
                                    <FormControlLabel key={value} value={value} control={<Radio />} label={value} />
                                )}
                            </RadioGroup>
                        </Box>
                        <Box>
                            <InputLabel sx={{ display: "inline" }}>API：</InputLabel>
                            <RadioGroup sx={{ display: "inline" }}
                                row
                                value={sysConfig.api}
                                onChange={e => setSysConfig(c => ({ ...c, api: e.target.value }))}
                                size="small"
                            >
                                {["mock", "web"].map(value =>
                                    <FormControlLabel key={value} value={value} control={<Radio />} label={value} />
                                )}
                            </RadioGroup>
                        </Box>
                        {process.env.NODE_ENV === 'production' ? [] :
                            <Box>
                                <InputLabel sx={{ display: "inline" }}>Web服务器：</InputLabel>
                                <TextField variant="standard" size="small"
                                    value={sysConfig.webServer}
                                    onChange={e => setSysConfig(c => ({ ...c, webServer: e.target.value }))}
                                />
                            </Box>
                        }
                        <Box>
                            <InputLabel sx={{ display: "inline" }}>客户端WebSocket端口：</InputLabel>
                            <TextField variant="standard" size="small"
                                value={sysConfig.port}
                                onChange={e => setSysConfig(c => ({ ...c, port: e.target.value }))}
                            />
                        </Box>
                    </>
                )}
        </Stack>
    );
};
