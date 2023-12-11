import { Box, Button, FormControlLabel, InputLabel, Radio, RadioGroup, Stack, TextField } from '@mui/material';
import useSysConfig from 'hox/sysConfig';

export default () => {
    const { sysConfig, setSysConfig, resetSysConfig } = useSysConfig();
    return (
        <Stack style={{ padding: "1em", overflowY: 'auto' }}>
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
                        {process.env.NODE_ENV === 'production' ? [] :
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
                        }
                        <Box>
                            <InputLabel sx={{ display: "inline" }}>Web服务器：</InputLabel>
                            <TextField variant="standard" size="small"
                                value={sysConfig.webServer}
                                onChange={e => setSysConfig(c => ({ ...c, webServer: e.target.value }))}
                            />
                        </Box>
                        {process.env.REACT_APP_PLATFORM === 'web' ? [] :
                            <>
                                <Box>
                                    <InputLabel sx={{ display: "inline" }}>Socket服务器：</InputLabel>
                                    <TextField variant="standard" size="small"
                                        value={sysConfig.socketServer}
                                        onChange={e => setSysConfig(c => ({ ...c, socketServer: e.target.value }))}
                                    />
                                </Box>
                                <Box>
                                    <InputLabel sx={{ display: "inline" }}>客户端Web服务器端口：</InputLabel>
                                    <TextField variant="standard" size="small"
                                        value={sysConfig.port}
                                        onChange={e => setSysConfig(c => ({ ...c, port: e.target.value }))}
                                    />
                                </Box>
                            </>
                        }
                    </>
                )}
        </Stack>
    );
};