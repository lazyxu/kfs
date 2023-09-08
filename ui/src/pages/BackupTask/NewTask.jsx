import { getSysConfig } from "hox/sysConfig";
import {
    Modal,
    Box,
    Button,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    Stack,
    TextField,
    DialogContent,
    Dialog,
    DialogTitle,
    IconButton,
    DialogActions
} from "@mui/material";
import { useEffect, useState } from "react";
import { getDriverApi } from "api/driver";
import AsyncSelect from "components/AsyncSelect";
import './index.scss';
import { newBackupTask } from "api/web/backup";
import { Close } from "@mui/icons-material";
import { noteError } from "components/Notification/Notification";

export default function ({ open, setOpen }) {
    const sysConfig = getSysConfig().sysConfig;
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [driverName, setDriverName] = useState('');
    const [concurrent, setConcurrent] = useState(2);
    const [encoder, setEncoder] = useState("none");
    const [srcPath, setSrcPath] = useState('');
    const [dstPath, setDstPath] = useState('/');
    return (
        <Dialog fullWidth={true} open={open} onClose={() => setOpen(false)} >
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>
                创建新的备份任务
                <IconButton
                    aria-label="close"
                    onClick={() => setOpen(false)}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close />
                </IconButton>
            </DialogTitle>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Stack spacing={2}>
                    <TextField variant="standard" label="备份名" type="search" sx={{ minWidth: "50%" }}
                        value={name}
                        onChange={e => setName(e.target.value)} />
                    <TextField variant="standard" label="描述" type="search" sx={{ minWidth: "50%" }}
                        value={description}
                        onChange={e => setDescription(e.target.value)} />
                    <TextField variant="standard" label="本地文件夹路径" type="search" sx={{ minWidth: "50%" }}
                        value={srcPath}
                        onChange={e => setSrcPath(e.target.value)} />
                    <Stack spacing={2} direction="row">
                        <FormControl sx={{ minWidth: "10em" }}>
                            <AsyncSelect
                                label="云盘"
                                fetchOptions={async () => {
                                    let drivers = await getDriverApi().listDriver();
                                    return drivers.map(driver => driver.name);
                                }}
                                value={driverName}
                                onChange={name => setDriverName(name)}
                            />
                        </FormControl>
                        <TextField variant="standard" label="远程文件夹路径" type="search" sx={{ minWidth: "50%" }}
                            value={dstPath}
                            onChange={e => setDstPath(e.target.value)}
                            disabled
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
                    <DialogActions sx={{
                        backgroundColor: theme => theme.background.primary,
                        color: theme => theme.context.primary
                    }}>
                        <Button variant="outlined" sx={{ width: "10em" }}
                            disabled={srcPath === "" || driverName === "" || name === ""}
                            onClick={e => {
                                newBackupTask(name, description, srcPath, driverName, dstPath, encoder, concurrent)
                                    .then(() => setOpen(false)).catch(e=>noteError(e.message))
                            }}
                        >
                            确定
                        </Button>
                    </DialogActions>
                </Stack>
            </DialogContent>
        </Dialog>
    )
}
