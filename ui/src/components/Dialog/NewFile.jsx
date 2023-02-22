import React, {useState} from 'react';
import {Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField} from "@mui/material";
import useDialog from "hox/dialog";
import {newFile} from "api/fs";
import useResourceManager from "hox/resourceManager";

export default () => {
    const [dialog, setDialog] = useDialog();
    let [name, setName] = useState("");
    const [resourceManager, setResourceManager] = useResourceManager();
    let {filePath, branchName} = resourceManager;
    return (
        <Dialog open={true} onClose={() => {
            setDialog(null)
        }}>
            <DialogTitle>{dialog.title}</DialogTitle>
            <DialogContent>
                <TextField
                    autoFocus={true}
                    margin="dense"
                    id="name"
                    placeholder="请输入文件的名字"
                    fullWidth
                    variant="standard"
                    onChange={e => {
                        setName(e.target.value)
                    }}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={() => {
                    setDialog(null);
                }}>取消</Button>
                <Button onClick={async () => {
                    await newFile(setResourceManager, branchName, filePath, name);
                    setDialog(null);
                }}>确定</Button>
            </DialogActions>
        </Dialog>
    );
};
