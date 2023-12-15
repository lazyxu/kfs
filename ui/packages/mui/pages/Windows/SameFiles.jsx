import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Stack } from "@mui/material";
import IconButton from "@mui/material/IconButton";

export default function ({ hash, sameFiles, onClose }) {
    return (
        <Dialog open={true} fullWidth={true} onClose={onClose}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>
                共 {sameFiles.length} 个相同文件
                <IconButton
                    aria-label="close"
                    onClick={() => onClose()}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close />
                </IconButton>
            </DialogTitle>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Stack sx={{ whiteSpace: "pre" }}  >
                    {sameFiles.map((f, i) => <Box key={i}>
                        {f.driverName}:{f.dirPath.length ? ("/" + f.dirPath.join("/") + "/" + f.name) : ("/" + f.name)}
                    </Box>)}
                </Stack>
            </DialogContent>
        </Dialog>
    );
}
