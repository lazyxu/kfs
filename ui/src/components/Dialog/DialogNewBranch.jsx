import React, {useCallback, useRef, useState} from 'react';
import {Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField} from "@mui/material";
import useDialog from "hox/dialog";
import {newDir} from "api/fs";
import useResourceManager from "hox/resourceManager";
import {newBranch} from "../../api/branch";

export default function () {
    const [dialog, setDialog] = useDialog();
    let [name, setName] = useState("");
    const [resourceManager, setResourceManager] = useResourceManager();
    let {filePath, branchName} = resourceManager;
    const inputRef = useRef(null);
    // useEffect(() => {
    //     console.log("inputRef.current", inputRef.current)
    //     inputRef.current?.focus();
    // });
    const autoFocusFn = useCallback(element => {
        console.log(element);
        window.xxx = element;
        element?.focus();
    }, []);
    return (
        <Dialog open={true} onClose={() => {
            setDialog(null)
        }}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>{dialog.title}</DialogTitle>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <TextField
                    // focused={true}
                    autoFocus
                    // inputProps={{autoFocus: true, focused: true}}
                    // autoFocus
                    // ref={inputRef}
                    margin="dense"
                    id="name"
                    placeholder="请输入文件夹的名字"
                    fullWidth
                    variant="outlined"
                    onChange={e => {
                        setName(e.target.value)
                    }}
                />
            </DialogContent>
            <DialogActions sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Button onClick={() => {
                    setDialog(null);
                }}>取消</Button>
                <Button onClick={() => {
                    setDialog(null);
                    newBranch(setResourceManager, name);
                }}>确定</Button>
            </DialogActions>
        </Dialog>
    );
};
