import useResourceManager from "@kfs/common/hox/resourceManager";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { Button, DialogActions, DialogContent, Link, Stack, Switch, TextField } from "@mui/material";
import { useState } from 'react';
import { newDriverBaiduPhoto } from "../../../api/driver";

// https://pan.baidu.com/union/doc/ol0rsap9s
export default function ({ onSucc }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { dirPath } = resourceManager;
    let [name, setName] = useState("新的同步文件夹");
    let [cron, setCron] = useState("*/10 * * * *");
    let [code, setCode] = useState("");
    let [syncAtStart, setSyncAtStart] = useState(true);
    const appKey = "huREKC2eNTctaBWfh3LdiAYjZ9ARBh5g";
    let redirectUri = `${getSysConfig().webServer}/api/v1/driver/baidu/callback`;
    redirectUri = `http://1zkl.com`;
    redirectUri = "oob";
    const url = `https://openapi.baidu.com/oauth/2.0/authorize?response_type=code&client_id=${appKey}&redirect_uri=${redirectUri}&scope=basic,netdisk&qrcode=1`;
    return (
        <>
            <DialogContent>
                <TextField
                    autoFocus
                    placeholder="文件夹名称"
                    margin="dense"
                    fullWidth
                    variant="outlined"
                    value={name}
                    onChange={e => setName(e.target.value)}
                />
                <Stack spacing={2} direction="row" sx={{ alignItems: "center" }}>
                    <TextField
                        placeholder="授权码"
                        variant="outlined"
                        onChange={e => setCode(e.target.value)}
                        sx={{ flex: 1 }}
                    />
                    {window.kfsEnv.VITE_APP_PLATFORM === 'web' ?
                        <Link target="_blank" href={url}>
                            点击获取授权码
                        </Link> :
                        <Link onClick={() => {
                            const { shell } = window.require('@electron/remote');
                            shell.openPath(url);
                        }}>点击获取授权码</Link>
                    }
                </Stack>
                <TextField
                    placeholder="cron"
                    margin="dense"
                    fullWidth
                    variant="outlined"
                    value={cron}
                    onChange={e => setCron(e.target.value)}
                />
                <Switch
                    checked={syncAtStart}
                    onChange={ (event) => {
                        setSyncAtStart(event.target.checked);
                    }}
                />
            </DialogContent>
            <DialogActions>
                <Button variant="outlined" sx={{ width: "10em" }} disabled={name === "" || code === ""} onClick={() => {
                    newDriverBaiduPhoto(dirPath, name, code, cron, syncAtStart).then(onSucc);
                }}>确定</Button>
            </DialogActions>
        </>
    )
}
