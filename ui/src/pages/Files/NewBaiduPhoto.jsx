import { Button, DialogActions, DialogContent, Link, TextField } from "@mui/material";
import { noteError } from 'components/Notification/Notification';
import useResourceManager from "hox/resourceManager";
import { useState } from 'react';
import { newDriver } from "../../api/driver";

// https://pan.baidu.com/union/doc/ol0rsap9s
export default function ({ setOpen }) {
    let [name, setName] = useState("");
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <TextField
                    autoFocus
                    margin="dense"
                    id="name"
                    placeholder="请输入云盘的名字"
                    fullWidth
                    variant="outlined"
                    onChange={e => {
                        setName(e.target.value)
                    }}
                />
                <Link href="https://openapi.baidu.com/oauth/2.0/authorize?response_type=code&client_id=iYCeC9g08h5vuP9UqvPHKKSVrKFXGa1v&redirect_uri=https://alist.nn.ci/tool/baidu/callback&scope=basic,netdisk&qrcode=1">
                    点击授权1
                </Link>
                <Button onClick={() => {
                    let redirectUri = "https://alist.nn.ci/tool/baidu/callback";
                    let clientId = "iYCeC9g08h5vuP9UqvPHKKSVrKFXGa1v";
                    let url = `https://openapi.baidu.com/oauth/2.0/authorize?response_type=code&client_id=${clientId}&redirect_uri=${redirectUri}&scope=basic,netdisk&qrcode=1`;
                    const { shell } = window.require('@electron/remote');
                    shell.openExternal(url);
                }} >
                    点击授权2
                </Button>
            </DialogContent>
            <DialogActions sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Button onClick={() => setOpen(false)}>取消</Button>
                <Button onClick={() => {
                    newDriver(setResourceManager, name)
                        .then(() => setOpen(false))
                        .catch(e => noteError("创建云盘失败：" + e.message));
                }}>确定</Button>
            </DialogActions>
        </>
    )
}
