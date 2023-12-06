import { Button, DialogActions, DialogContent, TextField } from "@mui/material";
import { useState } from 'react';
import { newDriver } from "../../../api/driver";

export default function ({ onSucc }) {
    let [name, setName] = useState("");
    let [description, setDescription] = useState("");
    return (
        <>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
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
            </DialogContent>
            <DialogActions sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Button variant="outlined" sx={{ width: "10em" }} disabled={name === ""} onClick={() => {
                    newDriver(name, description).then(onSucc);
                }}>确定</Button>
            </DialogActions>
        </>
    )
}
