import {Button, FormControl, InputLabel, MenuItem, Select, TextField} from "@mui/material";
import BackupSizeStatus from "./BackupSizeStatus";

export default function ({json}) {
    return (
        <>
            <TextField variant="standard" label="本地文件夹路径" type="search" sx={{width: "50%"}}
                       value={backupDir}
                       onChange={e => setBackupDir(e.target.value)}/>
            <Button variant="outlined" sx={{width: "10em"}}
                    disabled={isInvalidBackupDir(backupDir)}
                    onClick={e => {
                        sendJsonMessage({type: "calculateBackupSize.cancel", id, data: {backupDir: backupDir}});
                        const newId = id + 1;
                        setId(newId);
                        console.log("calculateBackupSize", newId, backupDir);
                        sendJsonMessage({type: "calculateBackupSize", id: newId, data: {backupDir: backupDir}});
                    }}
            >
                检测总大小
            </Button>
            <Button variant="outlined" sx={{width: "10em"}}
                    disabled={isInvalidBackupDir(backupDir)}
                    onClick={e => {
                        sendJsonMessage({type: "calculateBackupSize.cancel", id, data: {backupDir: backupDir}});
                    }}
            >
                取消
            </Button>
            <BackupSizeStatus json={lastJsonMessage}/>
            <FormControl sx={{width: "10em"}}>
                <InputLabel id="demo-simple-select-label">备份分支</InputLabel>
                <Select
                    labelId="demo-simple-select-label"
                    value={branchName}
                    onChange={e => setBranchName(e.target.value)}
                >
                    {branches.map(branch =>
                        <MenuItem key={branch.name} value={branch.name}>{branch.name}</MenuItem>
                    )}
                </Select>
            </FormControl>
        </>
    );
}
