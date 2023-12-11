import { Close } from '@mui/icons-material';
import { Box, ButtonGroup, Dialog, DialogContent, Stack } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import useWindows, { closeWindow } from 'hox/windows';


export function Window({ id, children }) {
    const [windows, setWindows] = useWindows();
    return (
        <Dialog open={true} fullScreen={true} onClose={() => closeWindow(setWindows, id)}>
            {children}
        </Dialog>
    )
}

export function TitleBar({ id, title, buttons }) {
    const [windows, setWindows] = useWindows();
    return (
        <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{
            color: theme => theme.context.secondary,
            backgroundColor: theme => theme.background.secondary,
        }}
        >
            <Box sx={{ paddingLeft: "1em" }}>
                {title}
            </Box>
            <Stack direction="row" justifyContent="flex-end" >
                <ButtonGroup variant="contained">
                    {buttons}
                </ButtonGroup>
                <IconButton aria-label="close" onClick={() => closeWindow(setWindows, id)}
                    sx={{
                        padding: "4px 12px", borderRadius: '0',
                        color: theme => theme.context.secondary,
                        ":hover": {
                            backgroundColor: "red",
                        }
                    }}
                >
                    <Close />
                </IconButton>
            </Stack>
        </Stack>
    )
}

export function WorkingArea({ children }) {
    return (
        <DialogContent sx={{
            padding: "0", paddingLeft: "5px",
            color: theme => theme.context.primary,
            backgroundColor: theme => theme.background.primary,
        }}>
            {children}
        </DialogContent>
    )
}

export function StatusBar({ children }) {
    return (
        <Box sx={{
            flex: "0 0 auto", padding: "8px",
            color: theme => theme.context.secondary,
            backgroundColor: theme => theme.background.secondary,
        }}>
            {children}
        </Box>
    )
}
