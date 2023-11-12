import { Button, DialogActions, DialogContent, Link, Stack, TextField } from "@mui/material";
import useResourceManager from "hox/resourceManager";
import { getSysConfig } from "hox/sysConfig";
import { useState } from 'react';
import { newDriverBaiduPhoto } from "../../api/driver";

// https://pan.baidu.com/union/doc/ol0rsap9s
export default function ({ setOpen }) {
    let [name, setName] = useState("");
    let [description, setDescription] = useState("");
    let [code, setCode] = useState("");
    const [resourceManager, setResourceManager] = useResourceManager();
    const appKey = "huREKC2eNTctaBWfh3LdiAYjZ9ARBh5g";
    let redirectUri = `${getSysConfig().sysConfig.webServer}/api/v1/driver/baidu/callback`;
    redirectUri = `http://1zkl.com`;
    redirectUri = "oob";
    const url = `https://openapi.baidu.com/oauth/2.0/authorize?response_type=code&client_id=${appKey}&redirect_uri=${redirectUri}&scope=basic,netdisk&qrcode=1`;
    return (
        <>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <TextField
                    autoFocus
                    placeholder="云盘名字"
                    margin="dense"
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
                <Stack spacing={2} direction="row" sx={{ alignItems: "center" }}>
                    <TextField
                        placeholder="授权码"
                        variant="outlined"
                        onChange={e => setCode(e.target.value)}
                        sx={{ flex: 1 }}
                    />
                    {process.env.REACT_APP_PLATFORM === 'web' ?
                        <Link target="_blank" href={url}>
                            点击获取授权码
                        </Link> :
                        <Link onClick={() => {
                            const { shell } = window.require('@electron/remote');
                            shell.openPath(url);
                        }}>点击获取授权码</Link>
                    }
                </Stack>
            </DialogContent>
            <DialogActions sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Button variant="outlined" sx={{ width: "10em" }} disabled={name === "" || code === ""} onClick={() => {
                    newDriverBaiduPhoto(setResourceManager, name, description, code).then(() => setOpen(false));
                }}>确定</Button>
            </DialogActions>
        </>
    )
}
