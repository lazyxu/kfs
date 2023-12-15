import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle } from "@mui/material";
import IconButton from "@mui/material/IconButton";

export default function ({ hash, metadata, onClose }) {
    return (
        <Dialog open={true} fullWidth={true} onClose={onClose}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>
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
                <Box
                    sx={{ whiteSpace: "pre" }}
                >
                    {JSON.stringify(metadata, null, 2)}
                </Box>
            </DialogContent>
        </Dialog>
    );
}
