import useWindows, { closeWindow } from '@kfs/mui/hox/windows';
import Close from '@mui/icons-material/Close';
import { Box, ButtonGroup, Dialog, DialogContent, Stack } from "@mui/material";
import IconButton from "@mui/material/IconButton";

export function Window({ wid, children }) {
    const [windows, setWindows] = useWindows();
    return (
        <Dialog open={true} fullScreen={true} onClose={() => closeWindow(setWindows, wid)}>
            {children}
        </Dialog>
    )
}

export function TitleBar({ wid, title, buttons }) {
    const [windows, setWindows] = useWindows();
    return (
        <Stack direction="row" justifyContent="space-between" alignItems="center">
            <Stack direction="row" justifyContent="flex-start" >
                <IconButton aria-label="close" onClick={() => closeWindow(setWindows, wid)}
                    sx={{
                        padding: "4px 12px", borderRadius: '0',
                        ":hover": {
                            backgroundColor: "red",
                        }
                    }}
                >
                    <Close />
                </IconButton>
                <ButtonGroup variant="contained">
                    {buttons}
                </ButtonGroup>
            </Stack>
            <Box sx={{ flex: 1, paddingLeft: "1em", }}>
                {title}
            </Box>
        </Stack>
    )
}

export function WorkingArea({ children }) {
    return (
        <DialogContent sx={{
            padding: "0", paddingLeft: "5px", backgroundColor: "var(--mui-palette-background-paper)"
        }}>
            {children}
        </DialogContent>
    )
}

export function StatusBar({ children }) {
    return (
        <Box sx={{
            flex: "0 0 auto", padding: "8px",
        }}>
            {children}
        </Box>
    )
}
