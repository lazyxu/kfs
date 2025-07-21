import useResourceManager from "@kfs/common/hox/resourceManager";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { Button, DialogActions, DialogContent, Link, Stack, Switch, TextField } from "@mui/material";
import { useState } from 'react';
import { createBaiduPhoto } from "../../../api/driver";

// https://pan.baidu.com/union/doc/ol0rsap9s
export default function ({ onClose }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { dirPath } = resourceManager;
    let [name, setName] = useState("来自：一刻相册");
    let [cron, setCron] = useState("0 4 * * *");
    let [bduss_bfess, setBdussBfess] = useState("");
    let [stoken, setStoken] = useState("");
    let [enabled, setEnabled] = useState(true);
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
                    placeholder="文件夹名字"
                    margin="dense"
                    fullWidth
                    variant="outlined"
                    value={name}
                    onChange={e => setName(e.target.value)}
                />
                <Link target="_blank" href="https://photo.baidu.com">一刻相册官网</Link>
                <TextField
                    autoFocus
                    placeholder="BDUSS_BFESS"
                    margin="dense"
                    fullWidth
                    variant="outlined"
                    value={bduss_bfess}
                    onChange={e => setBdussBfess(e.target.value)}
                />
                <TextField
                    autoFocus
                    placeholder="STOKEN"
                    margin="dense"
                    fullWidth
                    variant="outlined"
                    value={stoken}
                    onChange={e => setStoken(e.target.value)}
                />
                <TextField
                    placeholder="cron"
                    margin="dense"
                    fullWidth
                    variant="outlined"
                    value={cron}
                    onChange={e => setCron(e.target.value)}
                />
                <Switch
                    checked={enabled}
                    onChange={ (event) => {
                        setEnabled(event.target.checked);
                    }}
                />
                <Switch
                    checked={syncAtStart}
                    onChange={ (event) => {
                        setSyncAtStart(event.target.checked);
                    }}
                />
            </DialogContent>
            <DialogActions>
                <Button variant="outlined" sx={{ width: "10em" }} disabled={name === ""} onClick={() => {
                    createBaiduPhoto(dirPath, name, bduss_bfess, stoken, cron, enabled, syncAtStart).then(onClose);
                }}>确定</Button>
            </DialogActions>
        </>
    )
}
