import { Button, DialogActions, DialogContent, TextField } from "@mui/material";
import { noteError } from 'components/Notification/Notification';
import useResourceManager from "hox/resourceManager";
import { useState } from 'react';
import { newDriver } from "../../api/driver";

export default function ({ setOpen }) {
    let [name, setName] = useState("");
    let [description, setDescription] = useState("");
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
                <TextField
                    autoFocus
                    margin="dense"
                    id="name"
                    placeholder="请输入云盘的描述"
                    fullWidth
                    variant="outlined"
                    onChange={e => {
                        setDescription(e.target.value)
                    }}
                />
            </DialogContent>
            <DialogActions sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Button onClick={() => setOpen(false)}>取消</Button>
                <Button onClick={() => {
                    newDriver(setResourceManager, name, description)
                        .then(() => setOpen(false))
                        .catch(e => noteError("创建云盘失败：" + e.message));
                }}>确定</Button>
            </DialogActions>
        </>
    )
}
