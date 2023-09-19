import { Close } from '@mui/icons-material';
import { Dialog, DialogContent, DialogTitle, Grid } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import useResourceManager from "hox/resourceManager";
import moment from "moment/moment";

function formatTime(t) {
    return moment(t / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
}

function Attr({k, children}) {
    return <>
        <Grid xs={4} item sx={{overflowWrap: "anywhere"}}>{k}：</Grid>
        <Grid xs={8} item sx={{overflowWrap: "anywhere"}}>{children}</Grid>
    </>
}

export default ({setOpen, driver}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    // TODO: get more calculated attributes from server.
    return (
        <Dialog open={true} fullWidth={true} onClose={() => setOpen(false)}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>
                云盘属性
                <IconButton
                    aria-label="close"
                    onClick={() => setOpen(false)}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close/>
                </IconButton>
            </DialogTitle>
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Grid container spacing={1.5}>
                    <Attr k="云盘">{driver.name}</Attr>
                    <Attr k="描述">{driver.description}</Attr>
                </Grid>
            </DialogContent>
        </Dialog>
    );
};
