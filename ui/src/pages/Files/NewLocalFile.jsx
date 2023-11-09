import { FolderOpen } from "@mui/icons-material";
import { Button, DialogActions, DialogContent, FormControl, IconButton, InputLabel, MenuItem, Select, Stack, TextField } from "@mui/material";
import useResourceManager from "hox/resourceManager";
import { useState } from 'react';

export default function ({ setOpen }) {
    let [name, setName] = useState("");
    let [description, setDescription] = useState("");
    const [concurrent, setConcurrent] = useState(2);
    const [encoder, setEncoder] = useState("none");
    const [srcPath, setSrcPath] = useState('');
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Stack spacing={2}>
                    <TextField
                        autoFocus
                        margin="dense"
                        placeholder="云盘名字"
                        fullWidth
                        variant="outlined"
                        onChange={e => setName(e.target.value)}
                    />
                    <TextField
                        placeholder="云盘描述"
                        margin="dense"
                        fullWidth
                        variant="outlined"
                        onChange={e => setDescription(e.target.value)}
                    />
                    <Stack spacing={2} direction="row" sx={{}}>
                        <TextField variant="standard" label="本地文件夹路径" type="search" size="small" sx={{ flex: 1 }}
                            value={srcPath}
                            onChange={e => setSrcPath(e.target.value)} />
                        <IconButton component="label" variant="contained" onClick={async () => {
                            const { dialog } = window.require('@electron/remote');
                            const result = await dialog.showOpenDialog({
                                properties: ['openDirectory'],
                                defaultPath: srcPath,
                            });
                            if (!result.canceled) {
                                setSrcPath(result.filePaths[0]);
                            }
                        }}>
                            <FolderOpen />
                        </IconButton>
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
                </Stack>
            </DialogContent>
            <DialogActions sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Button variant="outlined" sx={{ width: "10em" }} disabled={srcPath === "" || name === ""} onClick={() => {
                    // newDriver(setResourceManager, name, description)
                        // .then(exist => exist ? noteError("云盘名称重复") : setOpen(false))
                    //     .catch(e => noteError("创建云盘失败：" + e.message));
                }}>确定</Button>
            </DialogActions>
        </>
    )
}
