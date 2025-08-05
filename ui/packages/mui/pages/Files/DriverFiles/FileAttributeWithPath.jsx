import { modeIsDir } from '@kfs/common/api/utils';
import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Tab, Tabs } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import { useState } from 'react';
import CommonFileAttribute from './components/FileAttribute';
import FolderCalculatedAttributes from './components/FolderCalculatedAttributes';
import FolderSyncAttributes from './components/FolderSyncAttributes';

function TabPanel({ children, value, index, ...other }) {
    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`file-tabpanel-${index}`}
            aria-labelledby={`file-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{ p: 0 }}>
                    {children}
                </Box>
            )}
        </div>
    );
}

export default ({ fileAttribute, onClose }) => {
    const { filePath, driverFile } = fileAttribute;
    const isDir = modeIsDir(driverFile.mode);
    const [tabValue, setTabValue] = useState(0);

    const handleTabChange = (event, newValue) => {
        setTabValue(newValue);
    };

    const tabs = [
        { label: '常规属性', component: <CommonFileAttribute fileAttribute={fileAttribute} /> }
    ];

    if (isDir) {
        tabs.push({ label: '计算属性', component: <FolderCalculatedAttributes driverFile={driverFile} filePath={filePath} /> });
    }

    if (driverFile.sync_from_id > 0) {
        tabs.push({ label: '同步设置', component: <FolderSyncAttributes driverFile={driverFile} onSave={onClose} /> });
    }

    return (
        <Dialog open={true} fullWidth={true} onClose={onClose} maxWidth="sm">
            <DialogTitle>
                属性
                <IconButton
                    aria-label="close"
                    onClick={() => onClose()}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                    }}
                >
                    <Close />
                </IconButton>
            </DialogTitle>
            <DialogContent>
                <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                    <Tabs value={tabValue} onChange={handleTabChange} aria-label="file attribute tabs">
                        {tabs.map((tab, index) => (
                            <Tab key={index} label={tab.label} id={`file-tab-${index}`} aria-controls={`file-tabpanel-${index}`} />
                        ))}
                    </Tabs>
                </Box>
                
                {tabs.map((tab, index) => (
                    <TabPanel key={index} value={tabValue} index={index}>
                        <Box sx={{ flexGrow: 1, mt: 2 }}>
                            {tab.component}
                        </Box>
                    </TabPanel>
                ))}
            </DialogContent>
        </Dialog>
    );
};